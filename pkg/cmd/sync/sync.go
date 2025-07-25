package sync

import (
	"fmt"
	"path/filepath"

	"github.com/jenkins-x-plugins/jx-gitops/pkg/helmfiles"
	"github.com/jenkins-x-plugins/jx-promote/pkg/environments"
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
	GitCommitUsername  string
	GitCommitUserEmail string
	AutoMerge          bool
	NoVersion          bool
	UpdateOnly         bool
	GitCredentials     bool
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
		Run: func(_ *cobra.Command, _ []string) {
			err := o.Run()
			helper.CheckErr(err)
		},
	}

	cmd.Flags().StringVar(&o.CommitTitle, "pull-request-title", "", "the PR title")
	cmd.Flags().StringVar(&o.CommitMessage, "pull-request-body", "", "the PR body")
	cmd.Flags().StringVarP(&o.GitCommitUsername, "git-user-name", "", "", "the user name to git commit")
	cmd.Flags().StringVarP(&o.GitCommitUserEmail, "git-user-email", "", "", "the user email to git commit")
	cmd.Flags().StringSliceVar(&o.Labels, "labels", []string{}, "a list of labels to apply to the PR")
	cmd.Flags().StringSliceVar(&o.ChartFilter.Namespaces, "namespaces", []string{}, "a list of namespaces to filter resources to sync")
	cmd.Flags().StringSliceVar(&o.ChartFilter.Charts, "charts", []string{}, "names of charts to filter resources to sync. Can be local chart name (without prefix) or the full name with prefix")
	cmd.Flags().BoolVarP(&o.AutoMerge, "auto-merge", "", true, "should we automatically merge if the PR pipeline is green")
	cmd.Flags().BoolVarP(&o.NoVersion, "no-version", "", false, "disables validation on requiring a '--version' option or environment variable to be required")
	cmd.Flags().BoolVarP(&o.UpdateOnly, "update-only", "", false, "only update versions in the target environment/namespace - do not add any new charts that are missing")
	cmd.Flags().BoolVarP(&o.GitCredentials, "git-credentials", "", false, "ensures the git credentials are setup so we can push to git")

	o.BaseOptions.AddBaseFlags(cmd)
	o.EnvironmentPullRequestOptions.ScmClientFactory.AddFlags(cmd)

	cmd.Flags().StringVarP(&o.CommitTitle, "commit-title", "", "", "the commit title")
	cmd.Flags().StringVarP(&o.CommitMessage, "commit-message", "", "", "the commit message")

	o.Source.AddFlags(cmd, "source")
	o.Target.AddFlags(cmd, "target")

	return cmd, o
}

// Validate validates the options
func (o *Options) Validate() error {
	err := o.BaseOptions.Validate()
	if err != nil {
		return fmt.Errorf("failed to validate base options: %w", err)
	}
	if o.Input == nil {
		o.Input = inputfactory.NewInput(&o.BaseOptions)
	}
	o.JXClient, o.Namespace, err = jxclient.LazyCreateJXClientAndNamespace(o.JXClient, o.Namespace)
	if err != nil {
		return fmt.Errorf("failed to create JX client: %w", err)
	}
	o.EnvMap, o.EnvNames, err = jxenv.GetOrderedEnvironments(o.JXClient, o.Namespace)
	if err != nil {
		return fmt.Errorf("failed to load environments: %w", err)
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
		return fmt.Errorf("failed to validate options: %w", err)
	}

	err = o.ChooseEnvironments()
	if err != nil {
		return fmt.Errorf("failed to choose environments: %w", err)
	}

	gitURL := o.Target.GitCloneURL
	if gitURL == "" {
		return fmt.Errorf("no target git clone URL")
	}

	o.SourceDir = ""
	sourceGitURL := o.Source.GitCloneURL
	if sourceGitURL != "" && gitURL != sourceGitURL {
		o.SourceDir, err = gitclient.CloneToDir(o.Git(), sourceGitURL, "")
		if err != nil {
			return fmt.Errorf("failed to clone source cluster %s: %w", sourceGitURL, err)
		}
	}

	// lets clear the branch name so we create a new one each time in a loop
	o.BranchName = ""

	if o.CommitTitle == "" {
		o.CommitTitle = "chore: sync versions"
	}

	o.Function = func() error {
		dir := o.OutDir
		return o.SyncVersions(o.SourceDir, dir)
	}

	_, err = o.EnvironmentPullRequestOptions.Create(gitURL, "", o.Labels, o.AutoMerge)
	if err != nil {
		return fmt.Errorf("failed to create Pull Request on repository %s: %w", gitURL, err)
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
		return fmt.Errorf("failed to gather target helmfiles from %s: %w", targetDir, err)
	}

	sourceHelmfiles := targetHelmfiles
	if sourceDir != "" {
		sourceHelmfiles, err = helmfiles.GatherHelmfiles(o.Source.Helmfile, sourceDir)
		if err != nil {
			return fmt.Errorf("failed to gather source helmfiles from %s: %w", sourceDir, err)
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
			return fmt.Errorf("failed to load repository prefixes from version stream dir %s: %w", o.VersionStreamDir, err)
		}
	}

	editor, err := helmfiles.NewEditor(targetDir, targetHelmfiles)
	if err != nil {
		return fmt.Errorf("failed to create helmfile editor: %w", err)
	}
	err = o.syncHelmfileVersions(sourceHelmfiles, editor)
	if err != nil {
		return fmt.Errorf("failed to sync versions: %w", err)
	}

	err = editor.Save()
	if err != nil {
		return fmt.Errorf("failed to save modified files: %w", err)
	}
	return nil
}

func (o *Options) syncHelmfileVersions(sourceHelmfiles []helmfiles.Helmfile, editor *helmfiles.Editor) error {
	for i := range sourceHelmfiles {
		src := &sourceHelmfiles[i]
		path := src.Filepath
		helmStates, err := helmfiles.LoadHelmfile(path)
		if err != nil {
			return fmt.Errorf("failed to load helmfile %s: %w", path, err)
		}

		for _, helmState := range helmStates {
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
				err = editor.AddChart(details)
				if err != nil {
					return fmt.Errorf("failed to add chart %s: %w", details.String(), err)
				}
			}
		}
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
