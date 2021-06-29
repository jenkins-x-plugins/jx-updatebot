package promote

import (
	"github.com/jenkins-x-plugins/jx-updatebot/pkg/argocd"
	"strings"

	"github.com/jenkins-x/jx-helpers/v3/pkg/kyamls"
	"github.com/jenkins-x/jx-logging/v3/pkg/log"
	"github.com/pkg/errors"
	"sigs.k8s.io/kustomize/kyaml/yaml"
)

func (o *Options) ModifyApplicationFiles(dir, repoURL string, version string) error {
	modifyFn := func(node *yaml.RNode, path string) (bool, error) {
		text := strings.TrimSpace(argocd.GetRepoURL(node, path))
		if argocd.TrimGitURLSuffix(repoURL) != argocd.TrimGitURLSuffix(text) {
			return false, nil
		}

		err := node.PipeE(yaml.LookupCreate(yaml.ScalarNode, "spec", "source", "targetRevision"), yaml.FieldSetter{StringValue: version})
		if err != nil {
			return false, errors.Wrapf(err, "failed to set spec.source.targetRevision to %s", version)
		}
		log.Logger().Infof("modified the version in file %s to %s", path, version)
		return true, nil
	}

	return kyamls.ModifyFiles(dir, modifyFn, argocd.ApplicationFilter)
}
