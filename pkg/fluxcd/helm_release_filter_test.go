package fluxcd_test

import (
	"testing"

	"github.com/jenkins-x-plugins/jx-updatebot/pkg/fluxcd"
	"github.com/jenkins-x-plugins/jx-updatebot/pkg/gitops"
	"github.com/stretchr/testify/assert"
)

func TestAppFilter(t *testing.T) {
	testCases := []struct {
		filter     fluxcd.HelmReleaseFilter
		matches    []fluxcd.ChartVersion
		notMatches []fluxcd.ChartVersion
	}{
		{
			filter: fluxcd.HelmReleaseFilter{},
			matches: []fluxcd.ChartVersion{
				{Chart: "https://github.com/myorg/myrepo", SourceRefName: "cheese"},
			},
		},
		{
			filter: fluxcd.HelmReleaseFilter{
				Chart: gitops.TextFilter{
					Includes: []string{"app1"},
				},
			},
			matches: []fluxcd.ChartVersion{
				{Chart: "https://github.com/myorg/app1", SourceRefName: "cheese"},
			},
			notMatches: []fluxcd.ChartVersion{
				{Chart: "https://github.com/myorg/app2", SourceRefName: "cheese"},
			},
		},
		{
			filter: fluxcd.HelmReleaseFilter{
				Chart: gitops.TextFilter{
					Excludes: []string{"app1"},
				},
			},
			matches: []fluxcd.ChartVersion{
				{Chart: "https://github.com/myorg/app2", SourceRefName: "cheese"},
			},
			notMatches: []fluxcd.ChartVersion{
				{Chart: "https://github.com/myorg/app1", SourceRefName: "cheese"},
			},
		},
	}

	for _, tc := range testCases {
		filter := tc.filter

		for _, v := range tc.matches {
			actual := filter.Matches(&v)
			assert.True(t, actual, "should match %#v", v)
		}

		for _, v := range tc.notMatches {
			actual := filter.Matches(&v)
			assert.False(t, actual, "should not match %#v", v)
		}
	}
}
