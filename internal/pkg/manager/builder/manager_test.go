package builder

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"

	"github.com/Grbisba/15-07-2025/internal/config"
	"github.com/Grbisba/15-07-2025/internal/model/dto"
)

func testManager(t *testing.T) *Manager {
	t.Helper()

	cfg := &config.Config{
		Loader: &config.Loader{
			MaxFilesPerTask:       3,
			MaxConcurrentTasks:    3,
			ZipDataFolder:         "/zip",
			ZipDataTeardown:       5,
			TempFilesDataFolder:   "/temp",
			TempFilesDataTeardown: 5,
			AppRootDir:            filepath.Join(os.Getenv("APP_ROOT_DIR"), "data"),
		}}
	dm, err := NewDirManager(zaptest.NewLogger(t), cfg)
	require.NoError(t, err)
	require.NotNil(t, dm)

	return dm
}

func TestManager_CreateFile(t *testing.T) {
	dm := testManager(t)
	var uuidString = uuid.New().String()

	t.Run("file", func(t *testing.T) {
		var s = "filedata"
		file := []byte(s)
		err := dm.CreateFile(file, uuidString, "test")
		assert.NoError(t, err)
		err = dm.CreateFile(file, uuidString, "test1")
		assert.NoError(t, err)
	})
	t.Run("zip", func(t *testing.T) {
		err := dm.CreateZipFile(uuidString)
		assert.NoError(t, err)
	})
}

func TestManager_BuildPathToZipFile(t *testing.T) {
	dm := testManager(t)
	require.NotNil(t, dm)

	t.Run("positive", func(t *testing.T) {
		var zipName = uuid.NewString()

		err := dm.CreateZipFile(zipName)
		assert.NoError(t, err)
		path, err := dm.buildPathToZipFile(zipName)
		assert.NoError(t, err)
		assert.Equal(t, path, filepath.Join(dm.zipDir, zipName+".zip"))
	})
}

func TestManager_BuildPathToTempFile(t *testing.T) {
	dm := testManager(t)

	t.Run("positive", func(t *testing.T) {
		var (
			fn  = uuid.NewString()
			dir = uuid.NewString()
		)

		err := dm.CreateFile([]byte(fn), dir, fn)
		assert.NoError(t, err)
		path, err := dm.buildPathToTempFile(dir, fn)
		assert.NoError(t, err)
		assert.Equal(t, path, filepath.Join(dm.fileDir, dir, fn))
	})
}

func TestManager_GetZipFile(t *testing.T) {
	dm := testManager(t)

	t.Run("positive", func(t *testing.T) {
		var zipName = uuid.NewString()

		err := dm.CreateZipFile(zipName)
		assert.NoError(t, err)
		err = dm.CreateFile([]byte("test"), zipName, "test")
		assert.NoError(t, err)
		err = dm.WriteToZipFile(zipName, dto.File{
			Name: "test",
			URL:  "123",
			Ext:  ".pdf",
		})
		assert.NoError(t, err)
		data, err := dm.GetZipFileData(zipName)
		assert.NoError(t, err)
		assert.NotEmpty(t, data)
	})
}

func TestManager_CreateZipFile(t *testing.T) {
	dm := testManager(t)

	t.Run("positive", func(t *testing.T) {
		var zipName = "b9f043cb-2fa6-424f-9ac8-107db430e30c"

		err := dm.CreateZipFile(zipName)
		assert.NoError(t, err)
	})
}
