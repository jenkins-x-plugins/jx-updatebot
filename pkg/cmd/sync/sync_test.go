package sync_test

import (
	"github.com/jenkins-x-plugins/jx-updatebot/pkg/cmd/sync"
	"github.com/jenkins-x/jx-helpers/v3/pkg/files"
	"github.com/jenkins-x/jx-helpers/v3/pkg/testhelpers"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

var (
	// generateTestOutput enable to regenerate the expected output
	generateTestOutput = true
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

		m := map[string]bool{}
		FindAllHelmfiles(t, m, outDir)
		FindAllHelmfiles(t, m, expectedDir)
		require.NotEmpty(t, m, "failed to find helmfile.yaml files")

		for f := range m {
			outFile := filepath.Join(outDir, f)
			expectedFile := filepath.Join(expectedDir, f)
			if generateTestOutput {
				dir := filepath.Dir(expectedFile)
				err = os.MkdirAll(dir, files.DefaultDirWritePermissions)
				require.NoError(t, err, "failed to make dir %s", dir)

				data, err := ioutil.ReadFile(outFile)
				require.NoError(t, err, "failed to load %s", outFile)

				err = ioutil.WriteFile(expectedFile, data, 0666)
				require.NoError(t, err, "failed to save file %s", expectedFile)
				t.Logf("saved %s\n", expectedFile)
			} else {
				t.Logf("verified %s\n", outFile)
				testhelpers.AssertEqualFileText(t, expectedFile, outFile)
			}
		}

		t.Logf("test %s has created %s\n", name, outDir)
	}
}

func FindAllHelmfiles(t *testing.T, m map[string]bool, dir string) {
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() || info.Name() != "helmfile.yaml" {
			return nil
		}
		rel, err := filepath.Rel(dir, path)
		if err != nil {
			return errors.Wrapf(err, "failed to get relative path of %s from %s", path, dir)
		}
		m[rel] = true
		return nil
	})
	require.NoError(t, err, "failed to walk dir %s", dir)

}
