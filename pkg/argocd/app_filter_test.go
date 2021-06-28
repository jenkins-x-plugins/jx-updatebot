package argocd_test

import (
	"github.com/jenkins-x-plugins/jx-updatebot/pkg/argocd"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAppFilter(t *testing.T) {
	testCases := []struct {
		filter     argocd.AppFilter
		matches    []argocd.AppVersion
		notMatches []argocd.AppVersion
	}{
		{
			filter: argocd.AppFilter{},
			matches: []argocd.AppVersion{
				{RepoURL: "https://github.com/myorg/myrepo", Path: "cheese"},
			},
		},
		{
			filter: argocd.AppFilter{
				RepoURL: argocd.TextFilter{
					Includes: []string{"app1"},
				},
			},
			matches: []argocd.AppVersion{
				{RepoURL: "https://github.com/myorg/app1", Path: "cheese"},
			},
			notMatches: []argocd.AppVersion{
				{RepoURL: "https://github.com/myorg/app2", Path: "cheese"},
			},
		},
		{
			filter: argocd.AppFilter{
				RepoURL: argocd.TextFilter{
					Excludes: []string{"app1"},
				},
			},
			matches: []argocd.AppVersion{
				{RepoURL: "https://github.com/myorg/app2", Path: "cheese"},
			},
			notMatches: []argocd.AppVersion{
				{RepoURL: "https://github.com/myorg/app1", Path: "cheese"},
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
