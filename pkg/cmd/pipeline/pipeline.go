package pipeline

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jenkins-x/lighthouse-client/pkg/triggerconfig/inrepo"

	"github.com/jenkins-x-plugins/jx-gitops/pkg/apis/gitops/v1alpha1"
	"github.com/jenkins-x-plugins/jx-gitops/pkg/plugins"
	"github.com/jenkins-x-plugins/jx-gitops/pkg/sourceconfigs"
	"github.com/jenkins-x-plugins/jx-pipeline/pkg/cmd/convert"
	"github.com/jenkins-x-plugins/jx-promote/pkg/environments"
	"github.com/jenkins-x/jx-helpers/v3/pkg/cmdrunner"
	"github.com/jenkins-x/jx-helpers/v3/pkg/cobras/helper"
	"github.com/jenkins-x/jx-helpers/v3/pkg/files"
	"github.com/jenkins-x/jx-helpers/v3/pkg/termcolor"
	"github.com/jenkins-x/jx-helpers/v3/pkg/yamls"
	"github.com/jenkins-x/jx-logging/v3/pkg/log"

	"github.com/spf13/cobra"
)

// Options the command line options
type Options struct {
	Dir        string
	ConfigFile string
	Filter     string
	KptBinary  string
	Strategy   string
	HomeDir    string
	AutoMerge  bool
	NoConvert  bool
	environments.EnvironmentPullRequestOptions
}

var info = termcolor.ColorInfo

// NewCmdUpgradePipeline creates a command object
func NewCmdUpgradePipeline() (*cobra.Command, *Options) {
	o := &Options{}

	cmd := &cobra.Command{
		Use:     "pipeline",
		Aliases: []string{"pipelines"},
		Short:   "Upgrades the pipelines in the source repositories to the latest version stream and pipeline catalog",
		Run: func(_ *cobra.Command, _ []string) {
			err := o.Run()
			helper.CheckErr(err)
		},
	}
	cmd.Flags().StringVarP(&o.Dir, "dir", "d", ".", "the directory look for the 'jx-requirements.yml` file")
	cmd.Flags().StringVarP(&o.Filter, "filter", "f", "", "the text filter to filter out repositories to upgrade")
	cmd.Flags().StringVarP(&o.ConfigFile, "config", "c", "", "the configuration file to load for the repository configurations. If not specified we look in .jx/gitops/source-repositories.yaml")
	cmd.Flags().StringVarP(&o.Strategy, "strategy", "s", "resource-merge", "the 'kpt' strategy to use. To see available strategies type 'kpt pkg update --help'. Typical values are: resource-merge, fast-forward, alpha-git-patch, force-delete-replace")

	cmd.Flags().StringVar(&o.CommitTitle, "pull-request-title", "", "the PR title")
	cmd.Flags().StringVar(&o.CommitMessage, "pull-request-body", "", "the PR body")
	cmd.Flags().BoolVarP(&o.AutoMerge, "auto-merge", "", true, "should we automatically merge if the PR pipeline is green")
	cmd.Flags().BoolVarP(&o.NoConvert, "no-convert", "", false, "disables converting from Kptfile based pipelines to the uses:sourceURI notation for reusing pipelines across repositories")
	cmd.Flags().StringVarP(&o.KptBinary, "bin", "", "", "the 'kpt' binary name to use. If not specified this command will download the jx binary plugin into ~/.jx3/plugins/bin and use that")

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

	config, err := o.LoadSourceConfig()
	if err != nil {
		return fmt.Errorf("failed to load source config: %w", err)
	}

	for i := range config.Spec.Groups {
		group := &config.Spec.Groups[i]
		for j := range group.Repositories {
			repo := &group.Repositories[j]

			if o.Filter != "" && !strings.Contains(repo.Name, o.Filter) {
				continue
			}
			err := o.UpgradeRepository(config, group, repo)
			if err != nil {
				log.Logger().Errorf("failed to upgrade repository %s due to: %s", repo.Name, err.Error())
			}
		}
	}
	return nil
}

func (o *Options) Validate() error {
	var err error
	if o.KptBinary == "" {
		o.KptBinary, err = plugins.GetKptBinary(plugins.KptVersion)
		if err != nil {
			return fmt.Errorf("failed to get kpt plugin: %w", err)
		}
	}
	if o.HomeDir == "" {
		o.HomeDir, err = os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get home dir: %w", err)
		}
	}

	// lazy create the git client
	o.EnvironmentPullRequestOptions.Git()

	if !o.NoConvert {
		defaultCatalog := "jenkins-x/jx3-pipeline-catalog"
		if inrepo.VersionStreamVersions[defaultCatalog] == "" {
			inrepo.VersionStreamVersions[defaultCatalog] = "HEAD"
		}
	}
	return nil
}

func (o *Options) LoadSourceConfig() (*v1alpha1.SourceConfig, error) {
	if o.ConfigFile == "" {
		o.ConfigFile = filepath.Join(o.Dir, ".jx", "gitops", v1alpha1.SourceConfigFileName)
	}

	exists, err := files.FileExists(o.ConfigFile)
	if err != nil {
		return nil, fmt.Errorf("failed to check if file exists %s: %w", o.ConfigFile, err)
	}

	config := &v1alpha1.SourceConfig{}
	if !exists {
		return nil, fmt.Errorf("no file %s please you sure you are running this command inside a git clone of your developent cluster repository", o.ConfigFile)
	}
	err = yamls.LoadFile(o.ConfigFile, config)
	if err != nil {
		return nil, fmt.Errorf("failed to load file %s: %w", o.ConfigFile, err)
	}
	return config, nil
}

func (o *Options) UpgradeRepository(config *v1alpha1.SourceConfig, group *v1alpha1.RepositoryGroup, repo *v1alpha1.Repository) error {
	err := sourceconfigs.DefaultValues(config, group, repo)
	if err != nil {
		return err
	}
	gitURL := repo.HTTPCloneURL
	if gitURL == "" {
		return nil
	}
	log.Logger().Infof("checking pipelines in repository: %s", info(gitURL))

	// lets clear the branch name so we create a new one each time in a loop
	o.BranchName = ""

	if o.CommitTitle == "" {
		o.CommitTitle = "chore: upgrade pipelines"
	}

	o.Function = func() error {
		dir := o.OutDir
		if o.NoConvert {
			return o.upgradePipelinesViaKpt(dir)
		}
		return o.convertPipelines(gitURL, dir)
	}

	_, err = o.EnvironmentPullRequestOptions.Create(gitURL, "", o.Labels, o.AutoMerge)
	if err != nil {
		return fmt.Errorf("failed to create Pull Request on repository %s: %w", gitURL, err)
	}
	return nil
}

func (o *Options) upgradePipelinesViaKpt(dir string) error {
	lhDir := filepath.Join(dir, ".lighthouse")
	exists, err := files.DirExists(lhDir)
	if err != nil {
		return fmt.Errorf("failed to check if dir %s exists: %w", lhDir, err)
	}
	if !exists {
		return nil
	}

	fs, err := os.ReadDir(lhDir)
	if err != nil {
		return fmt.Errorf("failed to read dir %s: %w", lhDir, err)
	}
	for _, f := range fs {
		if !f.IsDir() {
			continue
		}

		name := f.Name()
		kptFile := filepath.Join(lhDir, name, "Kptfile")
		exists, err := files.FileExists(kptFile)
		if err != nil {
			return fmt.Errorf("failed to check if file exists %s: %w", kptFile, err)
		}
		if !exists {
			continue
		}

		// clear the kpt repo cache everytime else we run into issues
		err = os.RemoveAll(filepath.Join(o.HomeDir, ".kpt", "repos"))
		if err != nil {
			return err
		}

		folder := filepath.Join(".lighthouse", name)

		args := []string{"pkg", "update", folder, "--strategy", o.Strategy}
		c := &cmdrunner.Command{
			Name: o.KptBinary,
			Args: args,
			Dir:  dir,
		}
		if o.CommandRunner == nil {
			o.CommandRunner = cmdrunner.DefaultCommandRunner
		}
		_, err = o.CommandRunner(c)
		if err != nil {
			return fmt.Errorf("failed to run %s: %w", c.CLI(), err)
		}
	}
	return nil
}

func (o *Options) convertPipelines(gitURL, dir string) error {
	lhDir := filepath.Join(dir, ".lighthouse")
	exists, err := files.DirExists(lhDir)
	if err != nil {
		return fmt.Errorf("failed to check if dir %s exists: %w", lhDir, err)
	}
	if !exists {
		return nil
	}

	_, co := convert.NewCmdPipelineConvertUses()

	co.Dir = dir

	err = co.Run()
	if err != nil {
		return fmt.Errorf("failed to update pipelines for repository %s in dir %s: %w", gitURL, dir, err)
	}
	return nil
}
