package argo

import (
	"github.com/jenkins-x-plugins/jx-updatebot/pkg/cmd/argo/promote"
	"github.com/jenkins-x-plugins/jx-updatebot/pkg/cmd/argo/sync"
	"github.com/jenkins-x/jx-helpers/v3/pkg/cobras"
	"github.com/jenkins-x/jx-logging/v3/pkg/log"
	"github.com/spf13/cobra"
)

// NewCmdArgo creates the new command
func NewCmdArgo() *cobra.Command {
	command := &cobra.Command{
		Use:     "argo",
		Aliases: []string{"argocd"},
		Short:   "Commands for working with ArgoCD",
		Run: func(command *cobra.Command, args []string) {
			err := command.Help()
			if err != nil {
				log.Logger().Errorf(err.Error())
			}
		},
	}
	command.AddCommand(cobras.SplitCommand(promote.NewCmdArgoPromote()))
	command.AddCommand(cobras.SplitCommand(sync.NewCmdArgoSync()))
	return command
}
