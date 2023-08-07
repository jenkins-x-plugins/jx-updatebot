package argocd

import (
	"github.com/jenkins-x/jx-helpers/v3/pkg/kyamls"
)

var (
	ApplicationFilter = kyamls.Filter{
		Kinds: []string{"argoproj.io/v1alpha1/Application", "argoproj.io/v1alpha1/ApplicationSet"},
	}
)
