package promote

import (
	"github.com/jenkins-x-plugins/jx-updatebot/pkg/fluxcd"
	"github.com/jenkins-x/jx-helpers/v3/pkg/kyamls"
	"sigs.k8s.io/kustomize/kyaml/yaml"
)

func (o *Options) ModifyApplicationFiles(dir, chart string, sourceRefName, version string) error {
	modifyFn := func(node *yaml.RNode, path string) (bool, error) {
		v := fluxcd.GetAppVersion(node, path)
		if chart != v.Chart {
			return false, nil
		}
		if sourceRefName != "" && sourceRefName != v.SourceRefName {
			return false, nil
		}
		err := fluxcd.SetAppVersion(node, path, version)
		if err != nil {
			return false, err
		}
		return true, nil
	}

	return kyamls.ModifyFiles(dir, modifyFn, fluxcd.HelmReleaseFilter)
}
