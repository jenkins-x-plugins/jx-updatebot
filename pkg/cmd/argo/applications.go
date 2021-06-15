package argo

import (
	"github.com/jenkins-x/jx-helpers/v3/pkg/kyamls"
	"github.com/jenkins-x/jx-logging/v3/pkg/log"
	"github.com/pkg/errors"
	"sigs.k8s.io/kustomize/kyaml/yaml"
	"strings"
)

var (
	filter = kyamls.Filter{
		Kinds: []string{"argoproj.io/v1alpha1/Application"},
	}
)

func (o *Options) ModifyApplicationFiles(dir, repoURL string, version string) error {
	modifyFn := func(node *yaml.RNode, path string) (bool, error) {

		value, err := node.Pipe(yaml.PathGetter{Path: []string{"spec", "source", "repoURL"}})
		if err != nil {
			return false, errors.Wrapf(err, "failed to get spec.source.repoURL")
		}
		text, err := value.String()
		if err != nil {
			return false, errors.Wrapf(err, "failed to get text value")
		}
		text = strings.TrimSpace(text)
		if trimGitURLSuffix(repoURL) != trimGitURLSuffix(text) {
			return false, nil
		}

		err = node.PipeE(yaml.LookupCreate(yaml.ScalarNode, "spec", "source", "targetRevision"), yaml.FieldSetter{StringValue: version})
		if err != nil {
			return false, errors.Wrapf(err, "failed to set spec.source.targetRevision to %s", version)
		}
		log.Logger().Infof("modified the version in file %s to %s", path, version)
		return true, nil
	}

	return kyamls.ModifyFiles(dir, modifyFn, filter)
}

// remove any trailing git tokens to make comparison less likely to fail
func trimGitURLSuffix(url string) string {
	return strings.TrimSuffix(strings.TrimSuffix(url, "/"), ".git")
}
