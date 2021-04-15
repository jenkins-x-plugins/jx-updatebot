package pr_test

import (
	"testing"

	"github.com/jenkins-x-plugins/jx-updatebot/pkg/cmd/pr"
	"github.com/jenkins-x/go-scm/scm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVersionTemplate(t *testing.T) {
	prNumber := 123
	owner := "myorg"
	repo := "my-repo"
	sha := "b054df5"
	fullName := scm.Join(owner, repo)
	prBranch := "my-pr-branch-name"
	expectedHeadClone := "https://github.com/" + fullName + ".git"

	pullRequest := &scm.PullRequest{
		Number: prNumber,
		Title:  "my awesome pull request",
		Body:   "some text",
		Source: prBranch,
		Head: scm.PullRequestBranch{
			Repo: scm.Repository{
				Clone:     expectedHeadClone,
				Namespace: owner,
				Name:      repo,
				FullName:  fullName,
			},
			Sha: sha,
		},
	}

	testCases := []struct {
		template string
		expected string
	}{
		{
			template: `{{ pullRequestSha "my-repo" }}`,
			expected: "b054df5",
		},
		{
			template: `{{ pullRequestSha "myorg/my-repo" }}`,
			expected: "b054df5",
		},
		{
			template: "1.2.3",
			expected: "1.2.3",
		},
	}

	for _, tc := range testCases {
		o := &pr.Options{}
		o.ScmClientFactory.NoWriteGitCredentialsFile = true

		o.AddPullRequest(pullRequest)

		actual, err := o.EvaluateVersionTemplate(tc.template, "sampleGitURL")
		require.NoError(t, err, "failed to evaluate template %s", tc.template)

		t.Logf("evaluated template '%s' and got '%s'\n", tc.template, actual)

		assert.Equal(t, tc.expected, actual, "for template %s", tc.template)
	}
}
