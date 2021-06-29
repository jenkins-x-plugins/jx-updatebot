package fluxcd_test

import (
	"github.com/jenkins-x-plugins/jx-updatebot/pkg/fluxcd"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAppFilter(t *testing.T) {
	testCases := []struct {
		filter     fluxcd.AppFilter
		matches    []fluxcd.AppVersion
		notMatches []fluxcd.AppVersion
	}{
		{
			filter: fluxcd.AppFilter{},
			matches: []fluxcd.AppVersion{
				{Chart: "https://github.com/myorg/myrepo", SourceRefName: "cheese"},
			},
		},
		{
			filter: fluxcd.AppFilter{
				Chart: fluxcd.TextFilter{
					Includes: []string{"app1"},
				},
			},
			matches: []fluxcd.AppVersion{
				{Chart: "https://github.com/myorg/app1", SourceRefName: "cheese"},
			},
			notMatches: []fluxcd.AppVersion{
				{Chart: "https://github.com/myorg/app2", SourceRefName: "cheese"},
			},
		},
		{
			filter: fluxcd.AppFilter{
				Chart: fluxcd.TextFilter{
					Excludes: []string{"app1"},
				},
			},
			matches: []fluxcd.AppVersion{
				{Chart: "https://github.com/myorg/app2", SourceRefName: "cheese"},
			},
			notMatches: []fluxcd.AppVersion{
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
