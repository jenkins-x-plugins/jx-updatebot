package cmd

import (
	"github.com/jenkins-x-plugins/jx-updatebot/pkg/cmd/pipeline"
	"github.com/jenkins-x-plugins/jx-updatebot/pkg/cmd/pr"
	"github.com/jenkins-x-plugins/jx-updatebot/pkg/cmd/version"
	"github.com/jenkins-x-plugins/jx-updatebot/pkg/rootcmd"
	"github.com/jenkins-x/jx-helpers/v3/pkg/cobras"
	"github.com/jenkins-x/jx-logging/v3/pkg/log"
	"github.com/spf13/cobra"
)

// Main creates the new command
func Main() *cobra.Command {
	cmd := &cobra.Command{
		Use:   rootcmd.TopLevelCommand,
		Short: "Updatebot commands",
		Run: func(cmd *cobra.Command, args []string) {
			err := cmd.Help()
			if err != nil {
				log.Logger().Errorf(err.Error())
			}
		},
	}
	cmd.AddCommand(cobras.SplitCommand(pipeline.NewCmdUpgradePipeline()))
	cmd.AddCommand(cobras.SplitCommand(pr.NewCmdPullRequest()))
	cmd.AddCommand(cobras.SplitCommand(version.NewCmdVersion()))
	return cmd
}
