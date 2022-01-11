package sync

import (
	"fmt"
	"path/filepath"

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
		Synchronizes some or all applications in an environment/namespace to another environment/namespace to reduce version drift

		Supports synchronizing environments or namespaces within the same cluster or namespaces between remote clusters (possibly using different namespaces).

		Create a Pull Request on the target GitOps repository to apply the changes so that you can review the changes before they happen. 
		You can use different labels to enable/disable auto-merging.
`)

	cmdExample = templates.Examples(`
		# choose the environments to synchronize
		jx updatebot sync

		# synchronizes the apps in 2 of your environments (local or remote)
		jx updatebot sync --source-env staging --target-env production

		# synchronizes the apps in 2 namespaces in the dev cluster
		jx updatebot sync --source-ns jx-staging --target-ns jx-production


		# synchronizes the edam and beer charts in 2 of your environments (local or remote)
		jx updatebot sync --source-env staging --target-env production --charts edam --charts beer

	`)
)

// Options the options for upgrading a cluster
type Options struct {
	options.BaseOptions
	environments.EnvironmentPullRequestOptions

	Source             EnvironmentOptions
	Target             EnvironmentOptions
	ChartFilter        ChartFilter
	PullRequestTitle   string
	PullRequestBody    string
	GitCommitUsername  string
	GitCommitUserEmail string
	AutoMerge          bool
	NoVersion          bool
	UpdateOnly         bool
	GitCredentials     bool
	Interactive        bool
	Labels             []string
	Input              input.Interface
	EnvMap             map[string]*v1.Environment
	EnvNames           []string
	SourceDir          string
	VersionStreamDir   string
	Prefixes           *versionstream.RepositoryPrefixes
}

type ChartFilter struct {
	Namespaces []string
	Charts     []string
}

// NewCmdEnvironmentSync creates a command object for the command
func NewCmdEnvironmentSync() (*cobra.Command, *Options) {
	o := &Options{}

	cmd := &cobra.Command{
		Use:     "sync",
		Short:   "Synchronizes some or all applications in an environment/namespace to another environment/namespace to reduce version drift",
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
	cmd.Flags().StringSliceVar(&o.ChartFilter.Namespaces, "namespaces", []string{}, "a list of namespaces to filter resources to sync")
	cmd.Flags().StringSliceVar(&o.ChartFilter.Charts, "charts", []string{}, "names of charts to filter resources to sync. Can be local chart name (without prefix) or the full name with prefix")
	cmd.Flags().BoolVarP(&o.AutoMerge, "auto-merge", "", true, "should we automatically merge if the PR pipeline is green")
	cmd.Flags().BoolVarP(&o.NoVersion, "no-version", "", false, "disables validation on requiring a '--version' option or environment variable to be required")
	cmd.Flags().BoolVarP(&o.UpdateOnly, "update-only", "", false, "only update versions in the target environment/namespace - do not add any new charts that are missing")
	cmd.Flags().BoolVarP(&o.GitCredentials, "git-credentials", "", false, "ensures the git credentials are setup so we can push to git")
	cmd.Flags().BoolVarP(&o.Interactive, "interactive", "", false, "enables interactive mode")
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

	// lazy create git
	o.EnvironmentPullRequestOptions.Git()
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
	if o.Source.Namespace != "" && stringhelpers.StringArrayIndex(o.ChartFilter.Namespaces, o.Source.Namespace) < 0 {
		o.ChartFilter.Namespaces = append(o.ChartFilter.Namespaces, o.Source.Namespace)
	}

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

	if o.VersionStreamDir == "" {
		o.VersionStreamDir = filepath.Join(sourceDir, "versionStream")
	}
	if o.Prefixes == nil {
		var err error
		o.Prefixes, err = versionstream.GetRepositoryPrefixes(o.VersionStreamDir)
		if err != nil {
			return errors.Wrapf(err, "failed to load repository prefixes from version stream dir %s", o.VersionStreamDir)
		}
	}

	editor, err := helmfiles.NewEditor(targetDir, targetHelmfiles)
	if err != nil {
		return errors.Wrapf(err, "failed to create helmfile editor")
	}
	err = o.syncHelmfileVersions(sourceHelmfiles, editor)
	if err != nil {
		return errors.Wrapf(err, "failed to sync versions")
	}

	err = editor.Save()
	if err != nil {
		return errors.Wrapf(err, "failed to save modified files")
	}
	return nil
}

func (o *Options) syncHelmfileVersions(sourceHelmfiles []helmfiles.Helmfile, editor *helmfiles.Editor) error {
	charts := make(map[string]*helmfiles.ChartDetails)
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
			if o.UpdateOnly {
				details.UpdateOnly = true
			}
			if !o.ChartFilter.Matches(details) {
				continue
			}
			if o.Target.Namespace != "" {
				details.Namespace = o.Target.Namespace
			}
			if !o.Interactive {
				err = editor.AddChart(details)
				if err != nil {
					return errors.Wrapf(err, "failed to add chart %s", details.String())
				}
			} else {
				charts[details.ReleaseName] = details
			}
		}
	}
	if o.Interactive {
		names := []string{}
		m := map[string]*helmfiles.ChartDetails{}
		for name, chart := range charts {
			text := chart.ReleaseName
			if chart.String() != "" {
				text = fmt.Sprintf("%-36s: %s", name, chart.Version)
			}
			names = append(names, text)
			m[text] = chart
		}
		results, err := o.Input.SelectNames(names, "Pick chart(s) to promote: ", false, "which chart name do you wish to promote")
		if err != nil {
			return err
		}
		for _, el := range results {
			err = editor.AddChart(m[el])
			if err != nil {
				return err
			}
		}
		return nil

	}
	return nil
}

// Matches return true if the chart details matches the filters
func (o *ChartFilter) Matches(details *helmfiles.ChartDetails) bool {
	if len(o.Namespaces) > 0 {
		if stringhelpers.StringArrayIndex(o.Namespaces, details.Namespace) < 0 {
			return false
		}
	}
	prefix, localName := helmfiles.SpitChartName(details.Chart)
	if len(o.Charts) > 0 {
		answer := false
		for _, c := range o.Charts {
			p, l := helmfiles.SpitChartName(c)
			if (prefix == p || p == "") && l == localName {
				answer = true
				break
			}
		}
		return answer
	}
	return true
}
