package flux

import (
	"github.com/jenkins-x-plugins/jx-updatebot/pkg/cmd/flux/promote"
	"github.com/jenkins-x-plugins/jx-updatebot/pkg/cmd/flux/sync"
	"github.com/jenkins-x/jx-helpers/v3/pkg/cobras"
	"github.com/jenkins-x/jx-logging/v3/pkg/log"
	"github.com/spf13/cobra"
)

// NewCmdFlux creates the new command
func NewCmdFlux() *cobra.Command {
	command := &cobra.Command{
		Use:     "flux",
		Aliases: []string{"fluxcd"},
		Short:   "Commands for working with FluxCD git repositories",
		Run: func(command *cobra.Command, _ []string) {
			err := command.Help()
			if err != nil {
				log.Logger().Errorf(err.Error())
			}
		},
	}
	command.AddCommand(cobras.SplitCommand(promote.NewCmdFluxPromote()))
	command.AddCommand(cobras.SplitCommand(sync.NewCmdFluxSync()))
	return command
}
