package pr

import (
	"context"
	"fmt"
	"github.com/jenkins-x-plugins/jx-updatebot/pkg/apis/updatebot/v1alpha1"
	"github.com/jenkins-x/jx-helpers/v3/pkg/cmdrunner"
	"github.com/jenkins-x/jx-helpers/v3/pkg/stringhelpers"
	"github.com/jenkins-x/jx-logging/v3/pkg/log"
	"github.com/pkg/errors"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
	"os"
	"strings"
)

// GoFindURLs find the git URLs for the given go dependency change
func (o *Options) GoFindURLs(rule *v1alpha1.Rule, change v1alpha1.Change, gc *v1alpha1.GoChange) error {
	ctx := context.Background()

	if o.GraphQLClient == nil {
		token := o.ScmClientFactory.GitToken
		if token == "" {
			token = os.Getenv("GIT_TOKEN")
		}
		if token == "" {
			token = os.Getenv("GITHUB_TOKEN")
		}
		ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
		hc := oauth2.NewClient(ctx, ts)
		o.GraphQLClient = githubv4.NewClient(hc)
	}

	for _, owner := range gc.Owners {
		if err := queryRepositoriesWithGoMod(ctx, o.GraphQLClient, rule, gc, owner); err != nil {
			return errors.Wrapf(err, "failed to query repositories")
		}
	}
	return nil
}

// ApplyGo applies the go change
func (o *Options) ApplyGo(dir string, gitURL string, change v1alpha1.Change, gc *v1alpha1.GoChange) error {
	o.CommitTitle = "chore: upgrade go dependencies"

	log.Logger().Infof("finding all the go dependences for repository: %s", gitURL)

	runner := cmdrunner.QuietCommandRunner
	c := &cmdrunner.Command{
		Dir:  dir,
		Name: "go",
		Args: []string{"list", "-m", "-f", "{{.Path}}", "all"},
	}
	text, err := runner(c)
	if err != nil {
		return errors.Wrapf(err, "failed to run command %s on %s", c.CLI(), gitURL)
	}

	lines := strings.Split(text, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && gc.UpgradePackages.Matches(line) {
			c = &cmdrunner.Command{
				Dir:  dir,
				Name: "go",
				Args: []string{"get", "-u=patch", line},
			}
			text, err = runner(c)
			if err != nil {
				log.Logger().Warnf("failed to update %s: %s", line, err.Error())
			}
			c = &cmdrunner.Command{
				Dir:  dir,
				Name: "go",
				Args: []string{"mod", "tidy"},
			}
			text, err = runner(c)
			if err != nil {
				log.Logger().Warnf("failed to update %s: %s", line, err.Error())
			}
		}
	}
	return nil
}

func queryRepositoriesWithGoMod(ctx context.Context, client *githubv4.Client, rule *v1alpha1.Rule, gc *v1alpha1.GoChange, owner string) error {
	var q struct {
		Organisation struct {
			Repositories struct {
				Edges []struct {
					Node struct {
						Name   string
						Object struct {
							Blob struct {
								Text string
							} `graphql:"... on Blob"`
						} `graphql:"object(expression: $fileFilter)"`
					}
				}
				PageInfo struct {
					EndCursor   githubv4.String
					HasNextPage bool
				}
			} `graphql:"repositories(first: 100, after: $commentsCursor)"`
		} `graphql:"organization(login: $owner)"`
	}
	v := map[string]interface{}{
		"owner":          githubv4.String(owner),
		"fileFilter":     githubv4.String("HEAD:go.mod"),
		"commentsCursor": (*githubv4.String)(nil), // Null after argument to get first page.
	}

	for {
		err := client.Query(ctx, &q, v)
		if err != nil {
			return errors.Wrapf(err, "github query failed")
		}

		for _, edge := range q.Organisation.Repositories.Edges {
			name := edge.Node.Name
			text := edge.Node.Object.Blob.Text
			if text == "" {
				continue
			}
			if !gc.Repositories.Matches(name) {
				continue
			}
			requirementsText := stripGoModuleLines(text)
			if strings.Contains(requirementsText, gc.Package) {
				log.Logger().Infof("about to process %s/%s", owner, name)

				u := fmt.Sprintf("https://github.com/%s/%s", owner, name)
				if stringhelpers.StringArrayIndex(rule.URLs, u) < 0 && stringhelpers.StringArrayIndex(rule.URLs, u+".git") < 0 {
					rule.URLs = append(rule.URLs, u)
				}
			}
		}

		if !q.Organisation.Repositories.PageInfo.HasNextPage {
			break
		}
		v["commentsCursor"] = githubv4.NewString(q.Organisation.Repositories.PageInfo.EndCursor)
	}
	return nil
}

func stripGoModuleLines(text string) string {
	buf := &strings.Builder{}
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "module ") {
			buf.WriteString(line)
			buf.WriteString("\n")
		}
	}
	return buf.String()
}
