package fluxcd

import (
	"github.com/jenkins-x/jx-helpers/v3/pkg/kyamls"
)

var (
	HelmReleaseFilter = kyamls.Filter{
		Kinds: []string{"helm.toolkit.fluxcd.io/v2beta1/HelmRelease"},
	}
)
