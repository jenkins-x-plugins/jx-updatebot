package sync

import (
	"fmt"

	"github.com/jenkins-x/jx-helpers/v3/pkg/options"
	"github.com/jenkins-x/jx-helpers/v3/pkg/stringhelpers"

	"github.com/spf13/cobra"
)

// EnvironmentOptions options to specify the environment git repository
type EnvironmentOptions struct {
	OptionPrefix    string
	GitCloneURL     string
	EnvironmentName string
	Namespace       string
	Dir             string
	Helmfile        string
}

// AddFlags adds the CLI flags to this object
func (o *EnvironmentOptions) AddFlags(cmd *cobra.Command, optionsPrefix string) {
	o.OptionPrefix = optionsPrefix
	cmd.Flags().StringVarP(&o.GitCloneURL, optionsPrefix+"-git-url", "", "", fmt.Sprintf("git URL to clone for the %s", optionsPrefix))
	cmd.Flags().StringVarP(&o.Helmfile, optionsPrefix+"-helmfile", "", "", "the helmfile to resolve. If not specified defaults to 'helmfile.yaml' in the git clone dir")
	cmd.Flags().StringVarP(&o.EnvironmentName, optionsPrefix+"-env", "", "", fmt.Sprintf("the environment name for the %s", optionsPrefix))
	cmd.Flags().StringVarP(&o.Dir, optionsPrefix+"-dir", "", "", fmt.Sprintf("the directory to use for the git clone for the %s", optionsPrefix))
	cmd.Flags().StringVarP(&o.Namespace, optionsPrefix+"-ns", "", "", fmt.Sprintf("the namespace for the %s", optionsPrefix))
}

func (o *EnvironmentOptions) IsBlank() bool {
	return o.EnvironmentName == "" && o.GitCloneURL == "" && o.Namespace == ""
}

func (o *Options) ChooseEnvironments() error {
	var err error
	if o.Source.IsBlank() {
		// lets pick a source environment
		o.Source.EnvironmentName, err = o.Input.PickNameWithDefault(o.EnvNames, "source environment: ", "", "pick the name of the source Environment you want to sync")
		if err != nil {
			return fmt.Errorf("failed to pick a source environment: %w", err)
		}
		if o.Source.EnvironmentName == "" {
			return fmt.Errorf("no source environment")
		}
	}
	if o.Target.IsBlank() {
		// lets pick a target environment
		targetEnvNames := o.EnvNames
		if o.Source.EnvironmentName != "" {
			targetEnvNames = stringhelpers.RemoveStringFromSlice(targetEnvNames, o.Source.EnvironmentName)
		}
		o.Target.EnvironmentName, err = o.Input.PickNameWithDefault(targetEnvNames, "target environment: ", "", "pick the name of the target Environment you want to sync")
		if err != nil {
			return fmt.Errorf("failed to pick a target environment: %w", err)
		}
		if o.Target.EnvironmentName == "" {
			return fmt.Errorf("no target environment")
		}
	}

	err = o.ValidateEnvironment(&o.Source, true)
	if err != nil {
		return fmt.Errorf("failed to validate the source: %w", err)
	}
	err = o.ValidateEnvironment(&o.Target, false)
	if err != nil {
		return fmt.Errorf("failed to validate the target: %w", err)
	}

	// lets validate the setup
	if o.Source.GitCloneURL == o.Target.GitCloneURL {
		if o.Source.Namespace == o.Target.Namespace {
			return fmt.Errorf("cannot use the same source and target git URL and namespace. You must sync with either different repositories or namespaces")
		}
	}
	return nil
}

// ValidateEnvironment lets validate we can find the helmfiles for the given
func (o *Options) ValidateEnvironment(env *EnvironmentOptions, source bool) error {
	name := "target"
	if source {
		name = "source"
	}
	var err error
	envName := env.EnvironmentName
	if env.GitCloneURL == "" && envName != "" {
		e := o.EnvMap[envName]
		if e == nil {
			return options.InvalidOption(name+"-env", envName, o.EnvNames)
		}

		if e.Spec.Namespace != "" {
			env.Namespace = e.Spec.Namespace
		}
		env.GitCloneURL = e.Spec.Source.URL
	}
	if o.Namespace == "" && envName == "" {
		return fmt.Errorf("no %s environment name or namespace for", name)
	}
	if env.GitCloneURL == "" {
		env.GitCloneURL, err = o.GetDevCloneGitURL()
		if err != nil {
			return err
		}
	}
	return nil
}

func (o *Options) GetDevCloneGitURL() (string, error) {
	e := o.EnvMap["dev"]
	if e == nil {
		return "", fmt.Errorf("no dev Environment found so cannot discover the cluster source git URL")
	}
	if e.Spec.Source.URL == "" {
		return "", fmt.Errorf("dev Environment has no cluster source git URL")
	}
	gitURL := e.Spec.Source.URL
	return gitURL, nil
}
