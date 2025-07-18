package builder

import (
	"os"
	"path/filepath"

	"go.uber.org/zap"
)

func (m *Manager) buildPathToTempFile(dirName, fn string) (string, error) {
	path := filepath.Join(m.fileDir, dirName, fn)

	info, err := os.Stat(path)
	if err != nil {
		m.log.Error(
			"error while stating temp file info by path",
			zap.String("path", path),
			zap.Error(err),
		)

		if os.IsNotExist(err) {
			return "", errNotExist
		}

		return "", errOsStat
	}

	if info != nil && info.Name() == fn && !info.IsDir() {
		return path, nil
	}

	m.log.Error(
		"unknown error while building temp file path",
		zap.String("path", path),
	)

	return "", errBuildPath
}

func (m *Manager) buildPathToZipFile(zipName string) (string, error) {
	path := filepath.Join(m.zipDir, zipName+".zip")

	info, err := os.Stat(path)
	if err != nil {
		m.log.Error(
			"error while stating zip file info by path",
			zap.String("path", path),
			zap.Error(err),
		)

		if os.IsNotExist(err) {
			return "", errNotExist
		}

		return "", errOsStat
	}

	if info != nil && info.Name() == zipName+".zip" && !info.IsDir() {
		return path, nil
	}

	m.log.Error(
		"unknown error while building zip file path",
		zap.String("path", path),
	)

	return "", errBuildPath
}

func (m *Manager) getDataFromFileByPath(path string) ([]byte, error) {
	fileData, err := os.ReadFile(path)
	if err != nil {
		m.log.Error(
			"failed to read file from path",
			zap.String("path", path),
			zap.Error(err),
		)

		return nil, errFileReading
	}

	return fileData, nil
}
