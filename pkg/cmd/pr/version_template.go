package pr

import (
	"github.com/Masterminds/sprig/v3"
	"github.com/jenkins-x/go-scm/scm"
	"github.com/jenkins-x/jx-helpers/v3/pkg/templater"
)

func (o *Options) EvaluateVersionTemplate(templateText, gitURL string) (string, error) {
	funcMap := sprig.TxtFuncMap()
	funcMap["pullRequestSha"] = func(name string) string {
		return o.PullRequestSHAs[name]
	}

	return templater.Evaluate(funcMap, o.TemplateData, templateText, "template.gotmpl", "version template for "+gitURL)
}

// AddPullRequest lets store pull requests so we can use the PR data later on
func (o *Options) AddPullRequest(pr *scm.PullRequest) {
	if o.PullRequestSHAs == nil {
		o.PullRequestSHAs = map[string]string{}
	}
	repo := pr.Repository()
	fullName := repo.FullName
	repoName := repo.Name
	if repoName == "" {
		repoName = pr.Head.Repo.Name
	}
	if fullName == "" {
		fullName = pr.Head.Repo.FullName
	}
	sha := pr.Head.Sha
	if sha != "" {
		o.PullRequestSHAs[repoName] = sha
		o.PullRequestSHAs[fullName] = sha
	}
}
