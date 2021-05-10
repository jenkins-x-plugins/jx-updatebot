package sync_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/jenkins-x-plugins/jx-gitops/pkg/helmfiles/testhelmfile"
	"github.com/jenkins-x-plugins/jx-updatebot/pkg/cmd/sync"
	"github.com/jenkins-x/jx-helpers/v3/pkg/files"
	"github.com/stretchr/testify/require"
)

var (
	// generateTestOutput enable to regenerate the expected output
	generateTestOutput = false
)

func TestSync(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "")
	require.NoError(t, err, "failed to create temp dir")

	testDir := "test_data"
	fileSlice, err := ioutil.ReadDir(testDir)
	require.NoError(t, err, "failed to read dir %s", testDir)

	testCaseName := os.Getenv("TEST_NAME")
	for _, f := range fileSlice {
		if !f.IsDir() {
			continue
		}
		name := f.Name()
		dir := filepath.Join(testDir, name)

		if testCaseName != "" && name != testCaseName {
			t.Logf("ignoring test case %s\n", name)
			continue
		}

		_, o := sync.NewCmdEnvironmentSync()

		switch name {
		case "name-filter":
			o.ChartFilter.Charts = []string{"ingress-nginx", "myapp"}
		case "ns-nginx":
			o.ChartFilter.Namespaces = []string{"nginx"}
		case "ns-staging":
			o.ChartFilter.Namespaces = []string{"jx-staging"}
		case "ns-prod":
			o.Source.Namespace = "jx-staging"
			o.Target.Namespace = "jx-production"
		case "update-only":
			o.UpdateOnly = true
			o.Source.Namespace = "jx-staging"
			o.Target.Namespace = "jx-production"
		}

		srcDir := filepath.Join(dir, "source")
		targetDir := filepath.Join(dir, "target")
		expectedDir := filepath.Join(dir, "expected")
		require.DirExists(t, srcDir)
		require.DirExists(t, targetDir)
		require.DirExists(t, expectedDir)

		outDir := filepath.Join(tmpDir, name)
		err := os.MkdirAll(outDir, files.DefaultDirWritePermissions)
		require.NoError(t, err, "failed to create dir %s", outDir)

		err = files.CopyDirOverwrite(targetDir, outDir)
		require.NoError(t, err, "failed to copy %s to %s", targetDir, outDir)

		err = o.SyncVersions(srcDir, outDir)
		require.NoError(t, err, "failed to process test %s", name)

		testhelmfile.AssertHelmfiles(t, expectedDir, outDir, generateTestOutput)
	}
}
