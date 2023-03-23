package promote_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jenkins-x-plugins/jx-updatebot/pkg/cmd/argo/promote"
	"github.com/jenkins-x/jx-helpers/v3/pkg/files"
	"github.com/jenkins-x/jx-helpers/v3/pkg/testhelpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestModifyApplicationFiles(t *testing.T) {
	tmpDir := t.TempDir()

	t.Logf("using dir %s\n", tmpDir)
	err := files.CopyDirOverwrite("test_data", tmpDir)
	require.NoError(t, err, "failed to copy test data to %s", tmpDir)

	dirNames, err := os.ReadDir(tmpDir)
	assert.NoError(t, err)

	repoURL := "https://github.com/myorg/myrepo.git"
	version := "v1.2.3"

	for _, d := range dirNames {
		if !d.IsDir() {
			continue
		}

		dir := d.Name()
		srcDir := filepath.Join(tmpDir, dir, "source")

		_, o := promote.NewCmdArgoPromote()

		err = o.ModifyApplicationFiles(srcDir, repoURL, version)
		require.NoError(t, err, "failed to modify files")

		fileNames, err := os.ReadDir(srcDir)
		require.NoError(t, err, "failed to read fileNames")

		for _, f := range fileNames {
			name := f.Name()
			if f.IsDir() || !strings.HasSuffix(name, ".yaml") {
				continue
			}
			expectedFile := filepath.Join(tmpDir, dir, "expected", name)
			srcFile := filepath.Join(srcDir, name)
			err = testhelpers.AssertEqualFileText(t, expectedFile, srcFile)
			require.NoError(t, err, "cannot assert expected file %s and actual file %s have the same text", expectedFile, srcFile)
		}
	}
}
