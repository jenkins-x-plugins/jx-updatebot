package pr_test

import (
	"testing"

	"github.com/jenkins-x-plugins/jx-updatebot/pkg/cmd/pr"
	"github.com/jenkins-x/go-scm/scm/driver/fake"
	"github.com/jenkins-x/jx-helpers/pkg/cmdrunner"
	"github.com/jenkins-x/jx-helpers/pkg/cmdrunner/fakerunner"
	"github.com/stretchr/testify/require"
)

func TestCreate(t *testing.T) {
	runner := &fakerunner.FakeRunner{
		CommandRunner: func(c *cmdrunner.Command) (string, error) {
			if c.Name == "git" && len(c.Args) > 0 && c.Args[0] == "push" {
				t.Logf("faking command %s in dir %s\n", c.CLI(), c.Dir)
				return "", nil
			}

			// lets really git clone but then fake out all other commands
			return cmdrunner.DefaultCommandRunner(c)
		},
	}

	_, o := pr.NewCmdPullRequest()
	o.Dir = "test_data"
	o.CommandRunner = runner.Run
	o.ScmClient, _ = fake.NewDefault()
	o.ScmClientFactory.ScmClient = o.ScmClient
	o.Version = "1.2.3"
	o.EnvironmentPullRequestOptions.ScmClientFactory.GitServerURL = "https://github.com"
	o.EnvironmentPullRequestOptions.ScmClientFactory.GitToken = "dummytoken"
	o.EnvironmentPullRequestOptions.ScmClientFactory.GitUsername = "dummyuser"

	err := o.Run()

	require.NoError(t, err, "failed to run command")
}
