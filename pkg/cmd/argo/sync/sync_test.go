package sync_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/jenkins-x-plugins/jx-updatebot/pkg/cmd/argo/sync"

	"github.com/jenkins-x/jx-helpers/v3/pkg/files"
	"github.com/stretchr/testify/require"
)

var (
	// generateTestOutput enable to regenerate the expected output
	generateTestOutput = false

	verbose = false
)

func TestArgoSync(t *testing.T) {
	tmpDir := t.TempDir()

	err := files.CopyDirOverwrite("test_data", tmpDir)
	require.NoError(t, err, "failed to copy test_data to %s", tmpDir)

	fileSlice, err := os.ReadDir(tmpDir)
	require.NoError(t, err, "failed to read dir %s", tmpDir)

	testCaseName := os.Getenv("TEST_NAME")
	for _, f := range fileSlice {
		if !f.IsDir() {
			continue
		}
		name := f.Name()
		dir := filepath.Join(tmpDir, name)

		if testCaseName != "" && name != testCaseName {
			t.Logf("ignoring test case %s\n", name)
			continue
		}

		_, o := sync.NewCmdArgoSync()

		switch name {
		case "include-app1":
			o.AppFilter.RepoURL.Includes = []string{"app1"}
		case "exclude-app1":
			o.AppFilter.RepoURL.Excludes = []string{"app1"}
		}

		srcDir := filepath.Join(dir, "source")
		targetDir := filepath.Join(dir, "target")
		expectedDir := filepath.Join("test_data", name, "expected")

		o.Source.Dir = srcDir
		o.Target.Dir = targetDir

		err = o.SyncVersions(srcDir, targetDir)
		require.NoError(t, err, "failed to run sync command")

		AssertDirContentsEqual(t, generateTestOutput, verbose, targetDir, expectedDir)
	}
}

// AssertDirContentsEqual asserts that the directory matches the expected dir
func AssertDirContentsEqual(t *testing.T, generateTestOutput, verbose bool, dir, expectedDir string) {
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error { //nolint:staticcheck
		if info == nil || info.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, ".yaml") && !strings.HasSuffix(path, ".yml") {
			return nil
		}

		rel, err := filepath.Rel(dir, path) //nolint:staticcheck
		require.NoError(t, err, "failed to return relative path")

		expectedFile := filepath.Join(expectedDir, rel)
		require.FileExists(t, path)

		resultData, err := os.ReadFile(path)
		result := strings.TrimSpace(string(resultData))
		require.NoError(t, err, "failed to load results %s", path)

		if generateTestOutput {
			expectedDir := filepath.Dir(expectedFile)
			err = os.MkdirAll(expectedDir, files.DefaultDirWritePermissions)
			require.NoError(t, err, "failed to create expected dir %s", expectedDir)

			err = os.WriteFile(expectedFile, []byte(result), 0600)
			require.NoError(t, err, "failed to save file %s", expectedFile)
			return nil
		}

		require.FileExists(t, expectedFile)
		expectData, err := os.ReadFile(expectedFile)
		require.NoError(t, err, "failed to load results %s", expectedFile)
		expectedText := strings.TrimSpace(string(expectData))

		if d := cmp.Diff(result, expectedText); d != "" {
			t.Errorf("modified file %s match expected: %s", path, d)
		}
		if verbose {
			t.Logf("found file %s file %s\n", path, result)
		}
		return nil
	})
	require.NoError(t, err, "failed to walk the source dir")
}
