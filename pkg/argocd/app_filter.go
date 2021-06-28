package argocd

import (
	"github.com/jenkins-x/jx-helpers/v3/pkg/stringhelpers"
	"github.com/spf13/cobra"
)

// TextFilter filters text
type TextFilter struct {
	Includes []string
	Excludes []string
}

// AppFilter filter for apps
type AppFilter struct {
	RepoURL TextFilter
	Path    TextFilter
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

func (o *TextFilter) AddFlags(cmd *cobra.Command, optionPrefix, name string) {
	cmd.Flags().StringSliceVar(&o.Includes, optionPrefix+"-include", []string{}, "text strings in the "+name+" to be included when synchronising")
	cmd.Flags().StringSliceVar(&o.Includes, optionPrefix+"-exclude", []string{}, "text strings in the "+name+" to be excluded when synchronising")
}

func (o *AppFilter) AddFlags(cmd *cobra.Command) {
	o.RepoURL.AddFlags(cmd, "repourl", "repository URL")
	o.RepoURL.AddFlags(cmd, "path", "path of the helm chart")
}
