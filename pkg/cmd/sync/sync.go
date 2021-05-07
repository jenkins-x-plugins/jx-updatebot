package sync

import (
	"fmt"
	"github.com/jenkins-x-plugins/jx-gitops/pkg/helmfiles"
	"github.com/jenkins-x-plugins/jx-promote/pkg/environments"
	"github.com/jenkins-x/go-scm/scm"
	v1 "github.com/jenkins-x/jx-api/v4/pkg/apis/jenkins.io/v1"
	"github.com/jenkins-x/jx-helpers/v3/pkg/cobras/helper"
	"github.com/jenkins-x/jx-helpers/v3/pkg/cobras/templates"
	"github.com/jenkins-x/jx-helpers/v3/pkg/gitclient"
	"github.com/jenkins-x/jx-helpers/v3/pkg/input"
	"github.com/jenkins-x/jx-helpers/v3/pkg/input/inputfactory"
	"github.com/jenkins-x/jx-helpers/v3/pkg/kube/jxclient"
	"github.com/jenkins-x/jx-helpers/v3/pkg/kube/jxenv"
	"github.com/jenkins-x/jx-helpers/v3/pkg/options"
	"github.com/jenkins-x/jx-helpers/v3/pkg/stringhelpers"
	"github.com/jenkins-x/jx-helpers/v3/pkg/versionstream"
	"github.com/jenkins-x/jx-helpers/v3/pkg/yaml2s"
	"github.com/pkg/errors"
	"github.com/roboll/helmfile/pkg/state"
	"github.com/spf13/cobra"
)

var (
	cmdLong = templates.LongDesc(`
		Synchonizes all or some applications in an environment from another environment to reduce version drift
`)

	cmdExample = templates.Examples(`
		# promotes your current cluster of Jenkins X to helm 3 / helmfile
		jx updatebot sync
	`)
)

// Options the options for upgrading a cluster
type Options struct {
	options.BaseOptions
	environments.EnvironmentPullRequestOptions

	Source             EnvironmentOptions
	Target             EnvironmentOptions
	PullRequestTitle   string
	PullRequestBody    string
	GitCommitUsername  string
	GitCommitUserEmail string
	AutoMerge          bool
	NoVersion          bool
	GitCredentials     bool
	Labels             []string
	Input              input.Interface
	EnvMap             map[string]*v1.Environment
	EnvNames           []string
	SourceDir          string
	Prefixes           *versionstream.RepositoryPrefixes
}

// NewCmdEnvironmentSync creates a command object for the command
func NewCmdEnvironmentSync() (*cobra.Command, *Options) {
	o := &Options{}

	cmd := &cobra.Command{
		Use:     "sync",
		Short:   "Synchonizes all or some applications in an environment from another environment",
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
	cmd.Flags().BoolVarP(&o.NoVersion, "no-version", "", false, "disables validation on requiring a '--version' option or environment variable to be required")
	cmd.Flags().BoolVarP(&o.GitCredentials, "git-credentials", "", false, "ensures the git credentials are setup so we can push to git")

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
	err := o.BaseOptions.Validate()
	if err != nil {
		return errors.Wrapf(err, "failed to validate base options")
	}
	if o.Input == nil {
		o.Input = inputfactory.NewInput(&o.BaseOptions)
	}
	o.JXClient, o.Namespace, err = jxclient.LazyCreateJXClientAndNamespace(o.JXClient, o.Namespace)
	if err != nil {
		return errors.Wrap(err, "failed to create JX client")
	}
	o.EnvMap, o.EnvNames, err = jxenv.GetOrderedEnvironments(o.JXClient, o.Namespace)
	if err != nil {
		return errors.Wrapf(err, "failed to load environments")
	}
	// lets remove the dev env name as we don't promote to/from it
	o.EnvNames = stringhelpers.RemoveStringFromSlice(o.EnvNames, "dev")
	return nil
}

// Run implements the command
func (o *Options) Run() error {
	err := o.Validate()
	if err != nil {
		return errors.Wrapf(err, "failed to validate options")
	}

	err = o.ChooseEnvironments()
	if err != nil {
		return errors.Wrapf(err, "failed to choose environments")
	}

	gitURL := o.Target.GitCloneURL
	if gitURL == "" {
		return errors.Errorf("no target git clone URL")
	}

	o.SourceDir = ""
	sourceGitURL := o.Source.GitCloneURL
	if sourceGitURL != "" && gitURL != sourceGitURL {
		o.SourceDir, err = gitclient.CloneToDir(o.Git(), sourceGitURL, "")
		if err != nil {
			return errors.Wrapf(err, "failed to clone source cluster %s", sourceGitURL)
		}
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
		return o.SyncVersions(o.SourceDir, dir)
	}

	_, err = o.EnvironmentPullRequestOptions.Create(gitURL, "", details, o.AutoMerge)
	if err != nil {
		return errors.Wrapf(err, "failed to create Pull Request on repository %s", gitURL)
	}
	return nil
}

// SyncVersions syncs the source and target versions
func (o *Options) SyncVersions(sourceDir, targetDir string) error {
	targetHelmfiles, err := helmfiles.GatherHelmfiles(o.Target.Helmfile, targetDir)
	if err != nil {
		return errors.Wrapf(err, "failed to gather target helmfiles from %s", targetDir)
	}

	sourceHelmfiles := targetHelmfiles
	if sourceDir != "" {
		sourceHelmfiles, err = helmfiles.GatherHelmfiles(o.Source.Helmfile, sourceDir)
		if err != nil {
			return errors.Wrapf(err, "failed to gather source helmfiles from %s", sourceDir)
		}
	} else {
		sourceDir = targetDir
	}

	editor, err := helmfiles.NewEditor(targetDir, targetHelmfiles)
	if err != nil {
		return errors.Wrapf(err, "failed to create helmfile editor")
	}
	err = o.syncHelmfileVersions(sourceDir, sourceHelmfiles, editor)
	if err != nil {
		return errors.Wrapf(err, "failed to sync versions")
	}

	err = editor.Save()
	if err != nil {
		return errors.Wrapf(err, "failed to save modified files")
	}
	return nil
}

func (o *Options) syncHelmfileVersions(sourceDir string, sourceHelmfiles []helmfiles.Helmfile, editor *helmfiles.Editor) error {
	for i := range sourceHelmfiles {
		src := &sourceHelmfiles[i]
		helmState := &state.HelmState{}
		path := src.Filepath
		err := yaml2s.LoadFile(path, helmState)
		if err != nil {
			return errors.Wrapf(err, "failed to load helmfile %s", path)
		}

		for i := range helmState.Releases {
			rel := &helmState.Releases[i]
			details := helmfiles.NewChartDetails(helmState, rel, o.Prefixes)
			if !o.Matches(details) {
				continue
			}
			err = editor.AddChart(details)
			if err != nil {
				return errors.Wrapf(err, "failed to add chart %s", details.String())
			}
		}
	}
	return nil
}

func (o *Options) Matches(details *helmfiles.ChartDetails) bool {
	// TODO add filters for the chart name / namespace
	return true
}
