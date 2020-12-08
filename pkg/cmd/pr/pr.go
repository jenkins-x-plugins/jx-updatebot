package pr

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/jenkins-x-plugins/jx-updatebot/pkg/apis/updatebot/v1alpha1"
	"github.com/jenkins-x-plugins/jx-updatebot/pkg/rootcmd"
	"github.com/jenkins-x/go-scm/scm"
	"github.com/jenkins-x/jx-helpers/v3/pkg/cobras/helper"
	"github.com/jenkins-x/jx-helpers/v3/pkg/cobras/templates"
	"github.com/jenkins-x/jx-helpers/v3/pkg/files"
	"github.com/jenkins-x/jx-helpers/v3/pkg/options"
	"github.com/jenkins-x/jx-helpers/v3/pkg/stringhelpers"
	"github.com/jenkins-x/jx-helpers/v3/pkg/termcolor"
	"github.com/jenkins-x/jx-helpers/v3/pkg/yamls"
	"github.com/jenkins-x/jx-logging/v3/pkg/log"
	"github.com/jenkins-x/jx-promote/pkg/environments"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/yargevad/filepathx"
)

var (
	info = termcolor.ColorInfo

	cmdLong = templates.LongDesc(`
		Create a Pull Request on each downstream repository
`)

	cmdExample = templates.Examples(`
		%s pr --test-url https://github.com/myorg/mytest.git
	`)
)

// Options the options for the command
type Options struct {
	environments.EnvironmentPullRequestOptions

	Dir              string
	ConfigFile       string
	Version          string
	VersionFile      string
	PullRequestTitle string
	PullRequestBody  string
	AutoMerge        bool
	Labels           []string
	TemplateData     map[string]interface{}
	PullRequestSHAs  map[string]string

	UpdateConfig v1alpha1.UpdateConfig
}

// NewCmdPullRequest creates a command object for the command
func NewCmdPullRequest() (*cobra.Command, *Options) {
	o := &Options{}

	cmd := &cobra.Command{
		Use:     "pr",
		Short:   "Create a Pull Request on each downstream repository",
		Long:    cmdLong,
		Example: fmt.Sprintf(cmdExample, rootcmd.BinaryName),
		Run: func(cmd *cobra.Command, args []string) {
			err := o.Run()
			helper.CheckErr(err)
		},
	}
	cmd.Flags().StringVarP(&o.Dir, "dir", "d", ".", "the directory look for the VERSION file")
	cmd.Flags().StringVarP(&o.ConfigFile, "config-file", "", "", "the updatebot config file. If none specified defaults to .jx/updatebot.yaml")
	cmd.Flags().StringVarP(&o.Version, "version", "", "", "the version number to promote. If not specified uses $VERSION or the version file")
	cmd.Flags().StringVarP(&o.VersionFile, "version-file", "", "", "the file to load the version from if not specified directly or via a $VERSION environment variable. Defaults to VERSION in the current dir")
	cmd.Flags().StringVar(&o.PullRequestTitle, "pull-request-title", "", "the PR title")
	cmd.Flags().StringVar(&o.PullRequestBody, "pull-request-body", "", "the PR body")
	cmd.Flags().StringSliceVar(&o.Labels, "labels", []string{}, "a list of labels to apply to the PR")
	cmd.Flags().BoolVarP(&o.AutoMerge, "auto-merge", "", true, "should we automatically merge if the PR pipeline is green")
	o.EnvironmentPullRequestOptions.ScmClientFactory.AddFlags(cmd)

	eo := &o.EnvironmentPullRequestOptions
	cmd.Flags().StringVarP(&eo.CommitTitle, "commit-title", "", "", "the commit title")
	cmd.Flags().StringVarP(&eo.CommitMessage, "commit-message", "", "", "the commit message")

	return cmd, o
}

// Run implements the command
func (o *Options) Run() error {
	err := o.Validate()
	if err != nil {
		return errors.Wrapf(err, "failed to validate")
	}

	for i, rule := range o.UpdateConfig.Spec.Rules {
		for _, gitURL := range rule.URLs {
			if gitURL == "" {
				log.Logger().Warnf("missing out repository %d as it has no git URL", i)
				continue
			}

			// lets clear the branch name so we create a new one each time in a loop
			o.BranchName = ""

			if o.PullRequestTitle == "" {
				o.PullRequestTitle = fmt.Sprintf("fix: upgrade to version %s", o.Version)
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

				for _, ch := range rule.Changes {
					err := o.ApplyChanges(dir, gitURL, ch)
					if err != nil {
						return errors.Wrapf(err, "failed to apply change")
					}

				}
				return nil
			}

			pr, err := o.EnvironmentPullRequestOptions.Create(gitURL, "", details, o.AutoMerge)
			if err != nil {
				return errors.Wrapf(err, "failed to create Pull Request on repository %s", gitURL)
			}
			if pr == nil {
				log.Logger().Infof("no Pull Request created")
				continue
			}
			o.AddPullRequest(pr)
		}
	}
	return nil
}

func (o *Options) Validate() error {
	if o.TemplateData == nil {
		o.TemplateData = map[string]interface{}{}
	}
	if o.PullRequestSHAs == nil {
		o.PullRequestSHAs = map[string]string{}
	}
	if o.Version == "" {
		if o.VersionFile == "" {
			o.VersionFile = filepath.Join(o.Dir, "VERSION")
		}
		exists, err := files.FileExists(o.VersionFile)
		if err != nil {
			return errors.Wrapf(err, "failed to check for file %s", o.VersionFile)
		}
		if exists {
			data, err := ioutil.ReadFile(o.VersionFile)
			if err != nil {
				return errors.Wrapf(err, "failed to read version file %s", o.VersionFile)
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

	// lets default the config file
	if o.ConfigFile == "" {
		o.ConfigFile = filepath.Join(o.Dir, ".jx", "updatebot.yaml")
	}
	exists, err := files.FileExists(o.ConfigFile)
	if err != nil {
		return errors.Wrapf(err, "failed to check for file %s", o.ConfigFile)
	}
	if exists {
		err = yamls.LoadFile(o.ConfigFile, &o.UpdateConfig)
		if err != nil {
			return errors.Wrapf(err, "failed to load config file %s", o.ConfigFile)
		}
	} else {
		log.Logger().Warnf("file %s does not exist so cannot create any updatebot Pull Requests", o.ConfigFile)
	}

	// lazy create the git client
	o.EnvironmentPullRequestOptions.Git()
	return nil
}

// ApplyChanges applies the changes to the given dir
func (o *Options) ApplyChanges(dir, gitURL string, change v1alpha1.Change) error {
	if change.Regex == nil {
		log.Logger().Infof("ignoring unknown change %#v", change)
		return nil
	}
	pattern := change.Regex.Pattern
	if pattern == "" {
		return errors.Errorf("no pattern for regex change %#v", change)
	}
	r, err := regexp.Compile(pattern)
	if err != nil {
		return errors.Wrapf(err, "failed to parse change regex: %s", pattern)
	}

	namedCaptures := make([]bool, 0)
	namedCapture := false
	for i, n := range r.SubexpNames() {
		if i == 0 {
			continue
		} else if n == "version" {
			namedCaptures = append(namedCaptures, true)
			namedCapture = true
		} else {
			namedCaptures = append(namedCaptures, false)
		}
	}

	for _, g := range change.Regex.Globs {
		path := filepath.Join(dir, g)
		matches, err := filepathx.Glob(path)
		if err != nil {
			return errors.Wrapf(err, "failed to evaluate glob %s", path)
		}
		for _, f := range matches {
			log.Logger().Infof("found file %s", f)

			data, err := ioutil.ReadFile(f)
			if err != nil {
				return errors.Wrapf(err, "failed to load file %s", f)
			}

			text := string(data)
			version := o.Version
			if change.VersionTemplate != "" {
				version, err = o.EvaluateVersionTemplate(change.VersionTemplate, gitURL)
				if err != nil {
					return errors.Wrapf(err, "failed to valuate version template %s", change.VersionTemplate)
				}
			}

			oldVersions := make([]string, 0)

			text2 := stringhelpers.ReplaceAllStringSubmatchFunc(r, text, func(groups []stringhelpers.Group) []string {
				answer := make([]string, 0)
				for i, group := range groups {
					if namedCapture {
						// If we are using named capture, then replace only the named captures that have the right name
						if namedCaptures[i] {
							oldVersions = append(oldVersions, group.Value)
							answer = append(answer, version)
						} else {
							answer = append(answer, group.Value)
						}
					} else {
						oldVersions = append(oldVersions, group.Value)
						answer = append(answer, version)
					}
				}
				return answer
			})

			if text2 != text {
				err = ioutil.WriteFile(f, []byte(text2), files.DefaultFileWritePermissions)
				if err != nil {
					return errors.Wrapf(err, "failed to save file %s", f)
				}
				log.Logger().Infof("modified file %s", info(f))
			}
		}
	}
	return nil
}
