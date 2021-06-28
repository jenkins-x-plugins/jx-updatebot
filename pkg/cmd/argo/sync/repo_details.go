package sync

import (
	"fmt"

	"github.com/spf13/cobra"
)

// RepositoryOptions options to specify the environment git repository
type RepositoryOptions struct {
	OptionPrefix string
	GitCloneURL  string
	Dir          string
}

// AddFlags adds the CLI flags to this object
func (o *RepositoryOptions) AddFlags(cmd *cobra.Command, optionsPrefix string) {
	o.OptionPrefix = optionsPrefix
	cmd.Flags().StringVarP(&o.GitCloneURL, optionsPrefix+"-git-url", "", "", fmt.Sprintf("git URL to clone for the %s", optionsPrefix))
	cmd.Flags().StringVarP(&o.Dir, optionsPrefix+"-dir", "", "", fmt.Sprintf("the directory to use for the git clone for the %s", optionsPrefix))
}
