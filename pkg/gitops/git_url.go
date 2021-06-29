package gitops

import "strings"

// remove any trailing git tokens to make comparison less likely to fail
func TrimGitURLSuffix(url string) string {
	return strings.TrimSuffix(strings.TrimSuffix(url, "/"), ".git")
}
