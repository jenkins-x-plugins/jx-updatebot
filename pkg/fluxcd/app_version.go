package fluxcd

import (
	"github.com/jenkins-x/jx-helpers/v3/pkg/kyamls"
	"github.com/jenkins-x/jx-logging/v3/pkg/log"
	"github.com/pkg/errors"
	"sigs.k8s.io/kustomize/kyaml/yaml"
)

// AppVersion represents an app version metadata from an Fluxcd Application
type AppVersion struct {
	Chart         string
	Version       string
	SourceRefName string
}

// Key returns a unique key for the app version
func (v *AppVersion) Key() string {
	return v.Chart + "\n" + v.SourceRefName
}

// String returns the string summary of the app version
func (v *AppVersion) String() string {
	sep := ""
	if v.SourceRefName != "" {
		sep = " sourceRefName: " + v.SourceRefName
	}
	return "repo: " + v.Chart + sep + " version: " + v.Version
}

// GetAppVersion gets the AppVersion from the given YAML file
func GetAppVersion(node *yaml.RNode, path string) *AppVersion {
	v := &AppVersion{}
	v.Chart = kyamls.GetStringField(node, path, "spec", "chart", "spec", "chart")
	v.Version = kyamls.GetStringField(node, path, "spec", "chart", "spec", "version")
	v.SourceRefName = kyamls.GetStringField(node, path, "spec", "chart", "spec", "sourceRef", "name")
	return v
}

// SetAppVersion sets the application version
func SetAppVersion(node *yaml.RNode, path, version string) error {
	err := node.PipeE(yaml.LookupCreate(yaml.ScalarNode, "spec", "chart", "spec", "version"), yaml.FieldSetter{StringValue: version})
	if err != nil {
		return errors.Wrapf(err, "failed to set spec.chart.spec.version to %s", version)
	}
	log.Logger().Debugf("modified the version in file %s to %s", path, version)
	return nil
}
