package v1alpha1

import (
	"github.com/jenkins-x/jx-helpers/v3/pkg/stringhelpers"
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
	// PullRequestLabels defines the labels to apply to created pull requests
	PullRequestLabels []string `json:"pullRequestLabels,omitempty"`

	// Rules defines the change rules
	Rules []Rule `json:"rules,omitempty"`
}

// Rule specifies a set of repositories and changes
type Rule struct {
	// URLs the git URLs of the repositories to create a Pull Request on
	URLs []string `json:"urls"`

	// Changes the changes to perform on the repositories
	Changes []Change `json:"changes"`

	// Fork if we should create the pull request from a fork of the repository
	Fork bool `json:"fork,omitempty"`

	// ReusePullRequest governs if existing pull requests for application are found and updated. Requires that --labels
	// or UpdateConfigSpec.PullRequestLabels are supplied.
	ReusePullRequest bool `json:"reusePullRequest,omitempty"`

	// SparseCheckout governs if sparse checkout is made of repository. Only possible with regex and go changes.
	// Note: Not all git servers support this.
	SparseCheckout bool `json:"sparseCheckout,omitempty"`

	// PullRequestAssignees
	PullRequestAssignees []string `json:"pullRequestAssignees,omitempty"`

	// AssignAuthorToPullRequests governs if downstream pull requests are automatically assigned to the upstream author
	AssignAuthorToPullRequests bool `json:"assignAuthorToPullRequests,omitempty"`
}

// Change the kind of change to make on a repository
type Change struct {
	// Command runs a shell command
	Command *Command `json:"command,omitempty"`

	// Go for go lang based dependency upgrades
	Go *GoChange `json:"go,omitempty"`

	// Regex a regex based modification
	Regex *Regex `json:"regex,omitempty"`

	// VersionStream updates the charts in a version stream repository
	VersionStream *VersionStreamChange `json:"versionStream,omitempty"`

	// VersionTemplate an optional template if the version is coming from a previous Pull Request SHA
	VersionTemplate string `json:"versionTemplate,omitempty"`
}

// Command runs a command line program
type Command struct {
	// Name the name of the command
	Name string `json:"name,omitempty"`
	// Args the command line arguments
	Args []string `json:"args,omitempty"`
	// Env the environment variables to pass into the command
	Env []EnvVar `json:"env,omitempty"`
}

// EnvVar the environment variable
type EnvVar struct {
	// Name the name of the environment variable
	Name string `json:"name,omitempty"`
	// Value the value of the environment variable
	Value string `json:"value,omitempty"`
}

// Regex a regex based modification
type Regex struct {
	// Pattern the regex pattern to apply
	Pattern string `json:"pattern,omitempty"`
	// Globs the files to apply this to
	Globs []string `json:"files,omitempty"`
}

// Pattern for matching strings
type Pattern struct {
	// Name
	Name string `json:"name,omitempty"`
	// Includes patterns to include in changing
	Includes []string `json:"include,omitempty"`
	// Excludes patterns to exclude from upgrading
	Excludes []string `json:"exclude,omitempty"`
}

// Matches returns true if the text matches the given text
func (p *Pattern) Matches(text string) bool {
	if p.Name != "" {
		return text == p.Name
	}
	return stringhelpers.StringMatchesAny(text, p.Includes, p.Excludes)
}

// VersionStreamChange for upgrading versions in a version stream
type VersionStreamChange struct {
	Pattern

	// Kind the kind of resources to change (charts, git, package etc)
	Kind string `json:"kind,omitempty"`
}

// GoChange for upgrading go dependencies
type GoChange struct {
	// Owners the git owners to query
	Owners []string `json:"owner,omitempty"`

	// Repositories the repositories to match
	Repositories Pattern `json:"repositories,omitempty"`

	// Package the text in the go.mod to filter on to perform an upgrade
	Package string `json:"package,omitempty"`

	// UpgradePackages the packages to upgrade
	UpgradePackages Pattern `json:"upgradePackages,omitempty"`

	// NoPatch disables patch upgrades so we can import to new minor releases
	NoPatch bool `json:"noPatch,omitempty"`
}
