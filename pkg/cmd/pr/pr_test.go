package pr_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jenkins-x/jx-helpers/v3/pkg/helmer"
	"github.com/jenkins-x/jx-helpers/v3/pkg/stringhelpers"

	"github.com/jenkins-x-plugins/jx-updatebot/pkg/cmd/pr"
	"github.com/jenkins-x/go-scm/scm/driver/fake"
	"github.com/jenkins-x/jx-helpers/v3/pkg/cmdrunner"
	"github.com/jenkins-x/jx-helpers/v3/pkg/cmdrunner/fakerunner"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreate(t *testing.T) {
	ev := os.Getenv("JX_EXCLUDE_TEST")
	if ev == "" {
		ev = "go"
	}
	excludeTests := strings.Split(ev, ",")
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

	fileNames, err := ioutil.ReadDir("test_data")
	assert.NoError(t, err)

	fakeHelmer := helmer.NewFakeHelmer()
	fakeHelmer.ChartsAllVersions["jxgh/jx-build-controller"] = []helmer.ChartSummary{
		{
			ChartVersion: "9.1.2",
		},
	}

	for _, f := range fileNames {
		if !f.IsDir() {
			continue
		}
		name := f.Name()
		if stringhelpers.StringArrayIndex(excludeTests, name) >= 0 {
			t.Logf("excluding test %s\n", name)
			continue
		}
		dir := filepath.Join("test_data", name)
		scmClient, fakeData := fake.NewDefault()

		_, o := pr.NewCmdPullRequest()
		o.Dir = dir
		o.CommandRunner = runner.Run
		o.ScmClient = scmClient
		o.ScmClientFactory.ScmClient = scmClient
		o.ScmClientFactory.NoWriteGitCredentialsFile = true
		o.Helmer = fakeHelmer
		o.Version = "1.2.3"
		o.EnvironmentPullRequestOptions.ScmClientFactory.GitServerURL = "https://github.com"
		o.EnvironmentPullRequestOptions.ScmClientFactory.GitToken = "dummytoken"
		o.EnvironmentPullRequestOptions.ScmClientFactory.GitUsername = "dummyuser"

		err := o.Run()
		require.NoError(t, err, "failed to run command for test %s", name)

		t.Logf("ran test %s\n", name)

		if name == "versionStream" {
			require.Len(t, fakeData.PullRequests, 1, "should have 1 Pull Request created for %s", name)
		}

		for n, pr := range fakeData.PullRequests {
			t.Logf("test %s created PR #%d with title: %s\n", name, n, pr.Title)
			t.Logf("body: %s\n\n", pr.Body)
		}

	}
}
