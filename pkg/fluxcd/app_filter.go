package fluxcd

import (
	"github.com/jenkins-x-plugins/jx-updatebot/pkg/gitops"
	"github.com/jenkins-x/jx-helpers/v3/pkg/stringhelpers"
	"github.com/spf13/cobra"
)

// AppFilter filter for apps
type AppFilter struct {
	Chart         gitops.TextFilter
	SourceRefName gitops.TextFilter
}

// Matches return true if the app version matches the filter
func (o *AppFilter) Matches(v *AppVersion) bool {
	if !stringhelpers.StringContainsAny(v.Chart, o.Chart.Includes, o.Chart.Excludes) {
		return false
	}
	if !stringhelpers.StringContainsAny(v.SourceRefName, o.SourceRefName.Includes, o.SourceRefName.Excludes) {
		return false
	}
	return true
}

func (o *AppFilter) AddFlags(cmd *cobra.Command) {
	o.Chart.AddFlags(cmd, "chart", "chart name")
	o.Chart.AddFlags(cmd, "source-ref-name", "the sourceRef name of the chart repository or bucket")
}
