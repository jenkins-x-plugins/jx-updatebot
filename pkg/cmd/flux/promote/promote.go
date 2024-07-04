package promote

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jenkins-x/jx-helpers/v3/pkg/cobras/templates"
	"github.com/jenkins-x/jx-helpers/v3/pkg/files"
	"github.com/jenkins-x/jx-logging/v3/pkg/log"

	"github.com/jenkins-x-plugins/jx-promote/pkg/environments"
	"github.com/jenkins-x/jx-helpers/v3/pkg/cobras/helper"
	"github.com/jenkins-x/jx-helpers/v3/pkg/kube/jxclient"
	"github.com/jenkins-x/jx-helpers/v3/pkg/options"

	"github.com/spf13/cobra"
)

// Options the command line options
type Options struct {
	Version       string
	VersionFile   string
	VersionPrefix string
	Dir           string
	Chart         string
	SourceRefName string
	TargetGitURL  string
	AutoMerge     bool
	environments.EnvironmentPullRequestOptions
}

var (
	cmdLong = templates.LongDesc(`
		Promotes a new HelmRelease version in a FluxCD git repository

		This command will use the given chart name and version along with an optional sourceRefName of the helm or git repository or bucket to find the HelmRelease resource in the target git repository and create a Pull Request if the version is different.
        This lets you push promotion pull requests into FluxCD repositories as part of your CI release pipeline.

		If you don't supply a version the $VERSION or VERSION file will be used. If you don't supply a chart the current folder name is used.
`)

	cmdExample = templates.Examples(`
		# lets promote a specific version of a chart with a source ref (repository) name to a git repo
		jx updatebot flux promote --version v1.2.3 --chart mychart --source-ref-name myrepo --target-git-url https://github.com/myorg/my-flux-repo.git

		# lets use the $VERSION env var or a VERSION file in the current dir and detect the chart name from the current folder
		jx updatebot flux promote --target-git-url https://github.com/myorg/my-flux-repo.git
	`)
)

// NewCmdFluxPromote creates a command object
func NewCmdFluxPromote() (*cobra.Command, *Options) {
	o := &Options{}

	cmd := &cobra.Command{
		Use:     "promote",
		Short:   "Promotes a new HelmRelease version in a FluxCD git repository",
		Long:    cmdLong,
		Example: cmdExample,
		Run: func(_ *cobra.Command, _ []string) {
			err := o.Run()
			helper.CheckErr(err)
		},
	}

	cmd.Flags().StringVarP(&o.Dir, "dir", "d", ".", "the directory look for the VERSION file")
	cmd.Flags().StringVarP(&o.Chart, "chart", "c", "", "the name of the chart to promote. If not specified defaults to the current directory name")
	cmd.Flags().StringVarP(&o.SourceRefName, "source-ref-name", "", "", "the source ref name of the HelmRepository, GitRepository or Bucket containing the helm chart")
	cmd.Flags().StringVarP(&o.TargetGitURL, "target-git-url", "", "", "the target git URL to create a Pull Request on")
	cmd.Flags().StringVarP(&o.Version, "version", "", "", "the version number to promote. If not specified uses $VERSION or the version file")
	cmd.Flags().StringVarP(&o.VersionFile, "version-file", "", "", "the file to load the version from if not specified directly or via a $VERSION environment variable. Defaults to VERSION in the current dir")
	cmd.Flags().StringVarP(&o.VersionPrefix, "version-prefix", "", "v", "the prefix added to the version number that will be used in the Flux CD Application YAML if --version option is not specified and the version is defaulted from $VERSION or the VERSION file")
	cmd.Flags().StringSliceVar(&o.Labels, "labels", []string{"promote"}, "a list of labels to apply to the PR")

	cmd.Flags().StringVar(&o.CommitTitle, "pull-request-title", "chore: upgrade the cluster git repository from the version stream", "the PR title")
	cmd.Flags().StringVar(&o.CommitMessage, "pull-request-body", "", "the PR body")
	cmd.Flags().BoolVarP(&o.AutoMerge, "auto-merge", "", false, "should we automatically merge if the PR pipeline is green")

	o.EnvironmentPullRequestOptions.ScmClientFactory.AddFlags(cmd)

	cmd.Flags().StringVarP(&o.CommitTitle, "commit-title", "", "", "the commit title")
	cmd.Flags().StringVarP(&o.CommitMessage, "commit-message", "", "", "the commit message")
	return cmd, o
}

// Run implements the command
func (o *Options) Run() error {
	err := o.Validate()
	if err != nil {
		return fmt.Errorf("failed to validate options: %w", err)
	}

	err = o.upgradeRepository(o.TargetGitURL)
	if err != nil {
		return fmt.Errorf("failed to create Pull Request on repository %s: %w", o.TargetGitURL, err)
	}
	return nil
}

func (o *Options) Validate() error {
	var err error

	if o.TargetGitURL == "" {
		return options.MissingOption("target-git-url")
	}
	if o.Chart == "" {
		if o.Dir == "" {
			o.Dir = "."
		}
		abs, err := filepath.Abs(o.Dir)
		if err != nil {
			return fmt.Errorf("failed to resolve absolute dir for %s: %w", o.Dir, err)
		}
		_, o.Chart = filepath.Split(abs)
	}
	if o.Chart == "" {
		return options.MissingOption("chart")
	}
	addPrefix := false
	if o.Version == "" {
		addPrefix = true
		if o.VersionFile == "" {
			o.VersionFile = filepath.Join(o.Dir, "VERSION")
		}
		exists, err := files.FileExists(o.VersionFile)
		if err != nil {
			return fmt.Errorf("failed to check for file %s: %w", o.VersionFile, err)
		}
		if exists {
			data, err := os.ReadFile(o.VersionFile)
			if err != nil {
				return fmt.Errorf("failed to read version file %s: %w", o.VersionFile, err)
			}
			o.Version = strings.TrimSpace(string(data))
		} else {
			log.Logger().Infof("version file %s does not exist", o.VersionFile)
		}
	}
	if o.Version == "" {
		o.Version = os.Getenv("VERSION")
		if o.Version == "" {
			return options.MissingOption("version")
		}
	}
	if addPrefix && o.VersionPrefix != "" && !strings.HasPrefix(o.Version, o.VersionPrefix) {
		o.Version = o.VersionPrefix + o.Version
	}

	o.EnvironmentPullRequestOptions.JXClient, o.EnvironmentPullRequestOptions.Namespace, err = jxclient.LazyCreateJXClientAndNamespace(o.EnvironmentPullRequestOptions.JXClient, o.EnvironmentPullRequestOptions.Namespace)
	if err != nil {
		return fmt.Errorf("failed to create jx client: %w", err)
	}

	// lazy create the git client
	o.EnvironmentPullRequestOptions.Git()
	return nil
}

func (o *Options) upgradeRepository(gitURL string) error {
	// lets clear the branch name so we create a new one each time in a loop
	o.BranchName = ""

	if o.CommitTitle == "" {
		o.CommitTitle = "chore: upgrade pipelines"
	}

	o.Function = func() error {
		dir := o.OutDir
		return o.ModifyHelmReleaseFiles(dir, o.Chart, o.SourceRefName, o.Version)
	}

	_, err := o.EnvironmentPullRequestOptions.Create(gitURL, "", o.Labels, o.AutoMerge)
	if err != nil {
		return fmt.Errorf("failed to create Pull Request on repository %s: %w", gitURL, err)
	}
	return nil
}
