package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// UpdateConfig defines the update rules
//
// +k8s:openapi-gen=true
type UpdateConfig struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata"`

	// Spec holds the update rule specifications
	Spec UpdateConfigSpec `json:"spec"`
}

// UpdateConfigSpec defines the rules to perform when updating.
type UpdateConfigSpec struct {

	// Rules defines the change rules
	Rules []Rule `json:"rules,omitempty"`
}

// Rule specifies a set of repositories and changes
type Rule struct {
	// URLs the git URLs of the repositories to create a Pull Request on
	URLs []string `json:"urls"`

	// Changes the changes to perform on the repositories
	Changes []Change `json:"changes"`
}

// Change the kind of change to make on a repository
type Change struct {
	// Regex a regex based modification
	Regex *Regex `json:"regex,omitempty"`

	// VersionStream updates the charts in a version stream repository
	VersionStream *VersionStreamChange `json:"versionStream,omitempty"`

	// VersionTemplate an optional template if the version is coming from a previous Pull Request SHA
	VersionTemplate string `json:"versionTemplate,omitempty"`
}

// Regex a regex based modification
type Regex struct {
	// Pattern the regex pattern to apply
	Pattern string `json:"pattern,omitempty"`
	// Globs the files to apply this to
	Globs []string `json:"files,omitempty"`
}

// VersionStreamChange for upgrading versions in a version stream
type VersionStreamChange struct {
	// Kind the kind of resources to change (charts, git, package etc)
	Kind string `json:"kind,omitempty"`
	// Name
	Name string `json:"name,omitempty"`
	// Includes patterns to include in changing
	Includes []string `json:"include,omitempty"`
	// Excludes patterns to exclude from upgrading
	Excludes []string `json:"exclude,omitempty"`
	// UpdateTektonImages if we should update the pipeline images too
	UpdateTektonImages bool `json:"updateTektonImages,omitempty"`
}
