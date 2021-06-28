package argocd

import (
	"github.com/jenkins-x/jx-helpers/v3/pkg/kyamls"
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
	return TrimGitURLSuffix(v.RepoURL) + "\n" + v.Path
}

// String returns the string summary of the app version
func (v *AppVersion) String() string {
	sep := ""
	if v.Path != "" {
		sep = " path: " + v.Path
	}
	return "repo: " + v.RepoURL + sep + " version: " + v.Version
}

// GetAppVersion gets the AppVersion from the given YAML file
func GetAppVersion(node *yaml.RNode, path string) *AppVersion {
	v := &AppVersion{}
	v.RepoURL = kyamls.GetStringField(node, path, "spec", "source", "repoURL")
	v.Path = kyamls.GetStringField(node, path, "spec", "source", "path")
	v.Version = kyamls.GetStringField(node, path, "spec", "source", "targetRevision")
	return v
}
