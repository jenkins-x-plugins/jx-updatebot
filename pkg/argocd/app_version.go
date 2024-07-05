package argocd

import (
	"fmt"

	"github.com/jenkins-x-plugins/jx-updatebot/pkg/gitops"
	"github.com/jenkins-x/jx-helpers/v3/pkg/kyamls"
	"github.com/jenkins-x/jx-logging/v3/pkg/log"

	"sigs.k8s.io/kustomize/kyaml/yaml"
)

// AppVersion represents an app version metadata from an ArgoCD Application
type AppVersion struct {
	RepoURL string
	Version string
	Path    string
}

// Key returns a unique key for the app version
func (v *AppVersion) Key() string {
	return gitops.TrimGitURLSuffix(v.RepoURL) + "\n" + v.Path
}

// String returns the string summary of the app version
func (v *AppVersion) String() string {
	sep := ""
	if v.Path != "" {
		sep = " path: " + v.Path
	}
	return "repo: " + v.RepoURL + sep + " version: " + v.Version
}

// GetRepoURL gets repository URL
func GetRepoURL(node *yaml.RNode, path string) string {

	annotation := kyamls.GetStringField(node, path, "metadata", "annotations", "gitops.jenkins-x.io/sourceRepoUrl")
	repoURL := kyamls.GetStringField(node, path, "spec", "source", "repoURL")
	if annotation != "" {
		return annotation
	}
	if repoURL != "" {
		return repoURL
	}
	return ""
}

// GetAppVersion gets the AppVersion from the given YAML file
func GetAppVersion(node *yaml.RNode, path string) *AppVersion {
	v := &AppVersion{}
	v.RepoURL = kyamls.GetStringField(node, path, "spec", "source", "repoURL")
	v.Path = kyamls.GetStringField(node, path, "spec", "source", "path")
	v.Version = kyamls.GetStringField(node, path, "spec", "source", "targetRevision")
	return v
}

// SetAppSetVersion sets the applicationSet version
func SetAppSetVersion(node *yaml.RNode, path, version string) error {
	err := node.PipeE(yaml.LookupCreate(yaml.ScalarNode, "spec", "template", "spec", "source", "targetRevision"), yaml.FieldSetter{StringValue: version})
	if err != nil {
		return fmt.Errorf("failed to set spec.generators.template.source.targetRevision to %s: %w", version, err)
	}
	log.Logger().Debugf("modified the version in file %s to %s", path, version)
	return nil
}

// SetAppVersion sets the application version
func SetAppVersion(node *yaml.RNode, path, version string) error {
	err := node.PipeE(yaml.LookupCreate(yaml.ScalarNode, "spec", "source", "targetRevision"), yaml.FieldSetter{StringValue: version})
	if err != nil {
		return fmt.Errorf("failed to set spec.source.targetRevision to %s: %w", version, err)
	}
	log.Logger().Debugf("modified the version in file %s to %s", path, version)
	return nil
}
