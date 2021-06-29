package sync

import (
	"fmt"
	"github.com/jenkins-x-plugins/jx-promote/pkg/environments"
	"github.com/jenkins-x-plugins/jx-updatebot/pkg/argocd"
	"github.com/jenkins-x-plugins/jx-updatebot/pkg/gitops"
	"github.com/jenkins-x/go-scm/scm"
	v1 "github.com/jenkins-x/jx-api/v4/pkg/apis/jenkins.io/v1"
	"github.com/jenkins-x/jx-helpers/v3/pkg/cobras/helper"
	"github.com/jenkins-x/jx-helpers/v3/pkg/cobras/templates"
	"github.com/jenkins-x/jx-helpers/v3/pkg/gitclient"
	"github.com/jenkins-x/jx-helpers/v3/pkg/input"
	"github.com/jenkins-x/jx-helpers/v3/pkg/input/inputfactory"
	"github.com/jenkins-x/jx-helpers/v3/pkg/kyamls"
	"github.com/jenkins-x/jx-helpers/v3/pkg/options"
	"github.com/jenkins-x/jx-helpers/v3/pkg/versionstream"
	"github.com/jenkins-x/jx-logging/v3/pkg/log"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"sigs.k8s.io/kustomize/kyaml/yaml"
)

var (
	cmdLong = templates.LongDesc(`
		Synchronizes some or all applications in an ArgoCD git repository to reduce version drift

		Creates a Pull Request on the target GitOps repository.
`)

	cmdExample = templates.Examples(`
		# create a Pull Request if any of the versions in the current directory are newer than the target repo
		jx updatebot argo sync --target-git-url https://github.com/myorg/my-production-repo

		# create a Pull Request if any of the versions are out of sync
		jx updatebot argo sync --source-git-url https://github.com/myorg/my-staging-repo --target-git-url https://github.com/myorg/my-production-repo

		# create a Pull Request if any of the versions are out of sync including only the given repo URL strings
		jx updatebot argo sync --source-git-url https://github.com/myorg/my-staging-repo --target-git-url https://github.com/myorg/my-production-repo --repourl-includes wine  --repourl-includes beer 

		# create a Pull Request if any of the versions are out of sync excluding the given repo URL strings
		jx updatebot argo sync --source-git-url https://github.com/myorg/my-staging-repo --target-git-url https://github.com/myorg/my-production-repo --repourl-excludes water
	`)
)

// Options the options for upgrading a cluster
type Options struct {
	options.BaseOptions
	environments.EnvironmentPullRequestOptions

	Source             gitops.RepositoryOptions
	Target             gitops.RepositoryOptions
	AppFilter          argocd.AppFilter
	PullRequestTitle   string
	PullRequestBody    string
	GitCommitUsername  string
	GitCommitUserEmail string
	AutoMerge          bool
	UpdateOnly         bool
	GitCredentials     bool
	Labels             []string
	Input              input.Interface
	EnvMap             map[string]*v1.Environment
	EnvNames           []string
	VersionStreamDir   string
	Prefixes           *versionstream.RepositoryPrefixes
	SourceApplications map[string]*argocd.AppVersion
}

// NewCmdArgoSync creates a command object for the command
func NewCmdArgoSync() (*cobra.Command, *Options) {
	o := &Options{}

	cmd := &cobra.Command{
		Use:     "sync",
		Short:   "Synchronizes some or all applications in an ArgoCD git repository to reduce version drift",
		Long:    cmdLong,
		Example: cmdExample,
		Run: func(cmd *cobra.Command, args []string) {
			err := o.Run()
			helper.CheckErr(err)
		},
	}

	cmd.Flags().StringVar(&o.PullRequestTitle, "pull-request-title", "", "the PR title")
	cmd.Flags().StringVar(&o.PullRequestBody, "pull-request-body", "", "the PR body")
	cmd.Flags().StringVarP(&o.GitCommitUsername, "git-user-name", "", "", "the user name to git commit")
	cmd.Flags().StringVarP(&o.GitCommitUserEmail, "git-user-email", "", "", "the user email to git commit")
	cmd.Flags().StringSliceVar(&o.Labels, "labels", []string{}, "a list of labels to apply to the PR")
	cmd.Flags().BoolVarP(&o.AutoMerge, "auto-merge", "", true, "should we automatically merge if the PR pipeline is green")
	// TODO support adding missing releases?
	//cmd.Flags().BoolVarP(&o.UpdateOnly, "update-only", "", false, "only update versions in the target environment/namespace - do not add any new charts that are missing")
	cmd.Flags().BoolVarP(&o.GitCredentials, "git-credentials", "", false, "ensures the git credentials are setup so we can push to git")

	o.AppFilter.AddFlags(cmd)

	o.BaseOptions.AddBaseFlags(cmd)
	o.EnvironmentPullRequestOptions.ScmClientFactory.AddFlags(cmd)

	eo := &o.EnvironmentPullRequestOptions
	cmd.Flags().StringVarP(&eo.CommitTitle, "commit-title", "", "", "the commit title")
	cmd.Flags().StringVarP(&eo.CommitMessage, "commit-message", "", "", "the commit message")

	o.Source.AddFlags(cmd, "source")
	o.Target.AddFlags(cmd, "target")

	return cmd, o
}

// Validate validates the options
func (o *Options) Validate() error {
	if o.SourceApplications == nil {
		o.SourceApplications = map[string]*argocd.AppVersion{}
	}
	err := o.BaseOptions.Validate()
	if err != nil {
		return errors.Wrapf(err, "failed to validate base options")
	}
	if o.Input == nil {
		o.Input = inputfactory.NewInput(&o.BaseOptions)
	}

	// lazy create git
	o.EnvironmentPullRequestOptions.Git()

	if o.Target.GitCloneURL == "" {
		return options.MissingOption(o.Target.OptionPrefix + "-git-url")
	}

	if o.Source.Dir == "" {
		sourceGitURL := o.Source.GitCloneURL
		if sourceGitURL == "" {
			// lets assume current directory is the source
			o.Source.Dir = "."
		} else {
			o.Source.Dir, err = gitclient.CloneToDir(o.Git(), sourceGitURL, "")
			if err != nil {
				return errors.Wrapf(err, "failed to clone source cluster %s", sourceGitURL)
			}
			if o.Source.Dir == "" {
				return errors.Errorf("failed to clone the source repository to a directory %s", sourceGitURL)
			}
		}
	}
	return nil
}

// Run implements the command
func (o *Options) Run() error {
	err := o.Validate()
	if err != nil {
		return errors.Wrapf(err, "failed to validate options")
	}

	gitURL := o.Target.GitCloneURL
	if gitURL == "" {
		return errors.Errorf("no target git clone URL")
	}

	// lets clear the branch name so we create a new one each time in a loop
	o.BranchName = ""

	if o.PullRequestTitle == "" {
		o.PullRequestTitle = fmt.Sprintf("chore: sync versions")
	}
	if o.CommitTitle == "" {
		o.CommitTitle = o.PullRequestTitle
	}
	source := ""
	details := &scm.PullRequest{
		Source: source,
		Title:  o.PullRequestTitle,
		Body:   o.PullRequestBody,
		Draft:  false,
	}

	for _, label := range o.Labels {
		details.Labels = append(details.Labels, &scm.Label{
			Name:        label,
			Description: label,
		})
	}

	o.Function = func() error {
		dir := o.OutDir
		return o.SyncVersions(o.Source.Dir, dir)
	}

	_, err = o.EnvironmentPullRequestOptions.Create(gitURL, "", details, o.AutoMerge)
	if err != nil {
		return errors.Wrapf(err, "failed to create Pull Request on repository %s", gitURL)
	}
	return nil
}

// SyncVersions syncs the source and target versions
func (o *Options) SyncVersions(sourceDir, targetDir string) error {
	err := o.findSourceApplications(sourceDir)
	if err != nil {
		return errors.Wrapf(err, "failed to find source Applications")
	}

	err = o.syncAppVersions(targetDir)
	if err != nil {
		return errors.Wrapf(err, "failed to modify target Applications")
	}
	return nil
}

func (o *Options) findSourceApplications(dir string) error {
	if o.SourceApplications == nil {
		o.SourceApplications = map[string]*argocd.AppVersion{}
	}
	modifyFn := func(node *yaml.RNode, path string) (bool, error) {
		v := argocd.GetAppVersion(node, path)
		if v.RepoURL == "" || v.Version == "" {
			return false, nil
		}

		log.Logger().Debugf("found source %s", v.String())

		k := v.Key()
		o.SourceApplications[k] = v
		return false, nil
	}
	return kyamls.ModifyFiles(dir, modifyFn, argocd.ApplicationFilter)
}

func (o *Options) syncAppVersions(dir string) error {
	modifyFn := func(node *yaml.RNode, path string) (bool, error) {
		v := argocd.GetAppVersion(node, path)
		if v.RepoURL == "" || !o.AppFilter.Matches(v) {
			return false, nil
		}
		k := v.Key()
		source := o.SourceApplications[k]
		if source == nil {
			return false, nil
		}

		err := argocd.SetAppVersion(node, path, source.Version)
		if err != nil {
			return false, err
		}
		return true, nil
	}
	return kyamls.ModifyFiles(dir, modifyFn, argocd.ApplicationFilter)
}
