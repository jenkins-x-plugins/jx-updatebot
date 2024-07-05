package pr

import (
	"fmt"
	"os"

	"github.com/jenkins-x-plugins/jx-updatebot/pkg/apis/updatebot/v1alpha1"
	"github.com/jenkins-x/jx-helpers/v3/pkg/cmdrunner"
)

func (o *Options) ApplyCommand(dir string, command *v1alpha1.Command) error {
	c := &cmdrunner.Command{
		Dir:  dir,
		Name: command.Name,
		Args: command.Args,
		Out:  os.Stdout,
		Err:  os.Stderr,
	}

	env := command.Env
	if len(env) > 0 {
		c.Env = map[string]string{}
		for _, e := range env {
			c.Env[e.Name] = e.Value
		}
	}

	_, err := o.CommandRunner(c)
	if err != nil {
		return fmt.Errorf("failed to run command %s: %w", c.CLI(), err)
	}
	return nil
}
