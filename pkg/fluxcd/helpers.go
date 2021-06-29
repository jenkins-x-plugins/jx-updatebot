package fluxcd

import (
	"github.com/jenkins-x/jx-helpers/v3/pkg/kyamls"
	"strings"
)

var (
	HelmReleaseFilter = kyamls.Filter{
		Kinds: []string{"helm.toolkit.fluxcd.io/v2beta1/HelmRelease"},
	}
)

// remove any trailing git tokens to make comparison less likely to fail
func TrimGitURLSuffix(url string) string {
	return strings.TrimSuffix(strings.TrimSuffix(url, "/"), ".git")
}
