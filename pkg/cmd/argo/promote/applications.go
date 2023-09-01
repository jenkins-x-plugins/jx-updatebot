package promote

import (
	"strings"

	"github.com/jenkins-x-plugins/jx-updatebot/pkg/argocd"
	"github.com/jenkins-x-plugins/jx-updatebot/pkg/gitops"

	"github.com/jenkins-x/jx-helpers/v3/pkg/kyamls"
	"sigs.k8s.io/kustomize/kyaml/yaml"
)

func (o *Options) ModifyApplicationFiles(dir, repoURL, version string) error {
	modifyFn := func(node *yaml.RNode, path string) (bool, error) {
		text := strings.TrimSpace(argocd.GetRepoURL(node, path))
		if gitops.TrimGitURLSuffix(repoURL) != gitops.TrimGitURLSuffix(text) {
			return false, nil
		}
		kind := kyamls.GetKind(node, path)
		var err error
		if kind == "Application" {
			err = argocd.SetAppVersion(node, path, version)
		} else if kind == "ApplicationSet" {
			err = argocd.SetAppSetVersion(node, path, version)
		}
		if err != nil {
			return false, err
		}
		return true, nil
	}

	return kyamls.ModifyFiles(dir, modifyFn, argocd.ApplicationFilter)
}
