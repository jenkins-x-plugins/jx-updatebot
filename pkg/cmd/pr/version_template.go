package pr

import (
	"github.com/Masterminds/sprig"
	"github.com/jenkins-x/go-scm/scm"
	"github.com/jenkins-x/jx-helpers/v3/pkg/maps"
	"github.com/jenkins-x/jx-helpers/v3/pkg/templater"
)

func (o *Options) EvaluateVersionTemplate(templateText, gitURL string) (string, error) {
	funcMap := sprig.TxtFuncMap()

	return templater.Evaluate(funcMap, o.TemplateData, templateText, "template.gotmpl", "version template for "+gitURL)
}

// AddPullRequest lets store pull requests so we can use the PR data later on
func (o *Options) AddPullRequest(pr *scm.PullRequest) {
	if o.TemplateData == nil {
		o.TemplateData = map[string]interface{}{}
	}
	repo := pr.Repository()
	repoName := repo.Name
	if repoName == "" {
		repoName = pr.Head.Repo.Name
	}
	sha := pr.Head.Sha
	if sha != "" {
		maps.SetMapValueViaPath(o.TemplateData, "PullRequests."+repoName+".Sha", sha)
	}
}
