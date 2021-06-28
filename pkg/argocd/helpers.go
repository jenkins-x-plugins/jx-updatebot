package argocd

import (
	"github.com/jenkins-x/jx-helpers/v3/pkg/kyamls"
	"strings"
)

var (
	ArgoApplicationFilter = kyamls.Filter{
		Kinds: []string{"argoproj.io/v1alpha1/Application"},
	}
)

// remove any trailing git tokens to make comparison less likely to fail
func TrimGitURLSuffix(url string) string {
	return strings.TrimSuffix(strings.TrimSuffix(url, "/"), ".git")
}
