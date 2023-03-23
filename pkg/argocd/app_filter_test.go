package argocd_test

import (
	"testing"

	"github.com/jenkins-x-plugins/jx-updatebot/pkg/argocd"
	"github.com/jenkins-x-plugins/jx-updatebot/pkg/gitops"
	"github.com/stretchr/testify/assert"
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
				RepoURL: gitops.TextFilter{
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
				RepoURL: gitops.TextFilter{
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

		for k := range tc.matches {
			v := tc.matches[k]
			actual := filter.Matches(&v)
			assert.True(t, actual, "should match %#v", v)
		}

		for k := range tc.notMatches {
			v := tc.notMatches[k]
			actual := filter.Matches(&v)
			assert.False(t, actual, "should not match %#v", v)
		}
	}
}
