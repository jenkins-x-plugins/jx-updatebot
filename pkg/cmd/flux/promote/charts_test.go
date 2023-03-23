package promote_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jenkins-x-plugins/jx-updatebot/pkg/cmd/flux/promote"
	"github.com/jenkins-x/jx-helpers/v3/pkg/files"
	"github.com/jenkins-x/jx-helpers/v3/pkg/testhelpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	// generateTestOutput enable to regenerate the expected output
	generateTestOutput = false
)

func TestModifyHelmReleaseFiles(t *testing.T) {
	tmpDir := t.TempDir()

	t.Logf("using dir %s\n", tmpDir)
	err := files.CopyDirOverwrite("test_data", tmpDir)
	require.NoError(t, err, "failed to copy test data to %s", tmpDir)

	dirNames, err := os.ReadDir(tmpDir)
	assert.NoError(t, err)

	chart := "chartmuseum"
	sourceRefName := ""
	version := "1.2.3"

	for _, d := range dirNames {
		if !d.IsDir() {
			continue
		}

		dir := d.Name()
		srcDir := filepath.Join(tmpDir, dir, "source")

		_, o := promote.NewCmdFluxPromote()

		err = o.ModifyHelmReleaseFiles(srcDir, chart, sourceRefName, version)
		require.NoError(t, err, "failed to modify files")

		fileNames, err := os.ReadDir(srcDir)
		require.NoError(t, err, "failed to read fileNames")

		for _, f := range fileNames {
			name := f.Name()
			if f.IsDir() || !strings.HasSuffix(name, ".yaml") {
				continue
			}
			expectedFile := filepath.Join("test_data", dir, "expected", name)
			srcFile := filepath.Join(srcDir, name)

			if generateTestOutput {
				data, err := os.ReadFile(srcFile)
				require.NoError(t, err, "failed to load %s", srcFile)

				err = os.WriteFile(expectedFile, data, 0600)
				require.NoError(t, err, "failed to save file %s", expectedFile)

				t.Logf("saved file %s\n", expectedFile)
				continue
			}

			err = testhelpers.AssertEqualFileText(t, expectedFile, srcFile)
			require.NoError(t, err, "cannot assert expected file %s and actual file %s have the same text", expectedFile, srcFile)
		}
	}
}
