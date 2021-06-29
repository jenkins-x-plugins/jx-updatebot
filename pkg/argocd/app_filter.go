package argocd

import (
	"github.com/jenkins-x-plugins/jx-updatebot/pkg/gitops"
	"github.com/jenkins-x/jx-helpers/v3/pkg/stringhelpers"
	"github.com/spf13/cobra"
)

// AppFilter filter for apps
type AppFilter struct {
	RepoURL gitops.TextFilter
	Path    gitops.TextFilter
}

// Matches return true if the app version matches the filter
func (o *AppFilter) Matches(v *AppVersion) bool {
	if !stringhelpers.StringContainsAny(v.RepoURL, o.RepoURL.Includes, o.RepoURL.Excludes) {
		return false
	}
	if !stringhelpers.StringContainsAny(v.Path, o.Path.Includes, o.Path.Excludes) {
		return false
	}
	return true
}

func (o *AppFilter) AddFlags(cmd *cobra.Command) {
	o.RepoURL.AddFlags(cmd, "repourl", "repository URL")
	o.RepoURL.AddFlags(cmd, "path", "path of the helm chart")
}
