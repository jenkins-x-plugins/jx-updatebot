package pr

import (
	"github.com/jenkins-x-plugins/jx-updatebot/pkg/apis/updatebot/v1alpha1"
	"github.com/jenkins-x/jx-helpers/v3/pkg/cmdrunner"
	"github.com/pkg/errors"
	"os"
)

func (o *Options) ApplyCommand(dir string, url string, change v1alpha1.Change, command *v1alpha1.Command) error {
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
		return errors.Wrapf(err, "failed to run command %s", c.CLI())
	}
	return nil
}
