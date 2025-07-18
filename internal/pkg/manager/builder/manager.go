package builder

import (
	"archive/zip"
	"bytes"
	"encoding/gob"
	"fmt"
	"os"
	"path/filepath"

	"go.uber.org/multierr"
	"go.uber.org/zap"

	"github.com/Grbisba/15-07-2025/internal/config"
	"github.com/Grbisba/15-07-2025/internal/model/dto"
	"github.com/Grbisba/15-07-2025/internal/pkg/manager"
)

var _ manager.Manager = (*Manager)(nil)

type Manager struct {
	rootDir string
	cfg     *config.Loader
	log     *zap.Logger
	zipDir  string
	fileDir string
}

func NewDirManager(log *zap.Logger, cfg *config.Config) (*Manager, error) {
	if cfg == nil {
		return nil, errNilReference
	}
	if log == nil {
		log = zap.NewNop()
	}

	m := &Manager{
		rootDir: cfg.Loader.AppRootDir,
		cfg:     cfg.Loader,
		log:     log.Named("dir-manager"),
	}

	var err error
	err = os.MkdirAll(m.cfg.AppRootDir, os.ModePerm)
	if err != nil {
		return nil, err
	}

	m.zipDir, err = m.createTempZipDir()
	if err != nil {
		return nil, err
	}

	m.fileDir, err = m.createTempFileDir()
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (m *Manager) CreateFile(fileData []byte, dirName string, fn string) error {
	dirPath := filepath.Join(m.fileDir, dirName)
	err := m.createDirByPath(dirPath)
	if err != nil {
		return err
	}

	filePath := filepath.Join(dirPath, fn)

	m.log.Debug("creating file", zap.String("dir", dirName), zap.String("name", fn))
	out, err := os.Create(filePath)
	if err != nil {
		m.log.Error(
			"failed to create zip file",
			zap.String("dir", dirName),
			zap.String("file", fn),
			zap.Error(err),
		)

		return err
	}

	defer func() {
		err = multierr.Combine(err, out.Close())
	}()

	err = gob.NewEncoder(out).Encode(fileData)
	if err != nil {
		m.log.Error(
			"failed to encode file data",
			zap.String("dir", dirName),
			zap.String("file", fn),
			zap.Error(err),
		)

		return errGobEncDec
	}

	return err
}

func (m *Manager) CreateZipFile(zipName string) error {
	zipPath := filepath.Join(m.zipDir, fmt.Sprintf("%s.zip", zipName))
	outFile, err := os.Create(zipPath)
	if err != nil {
		m.log.Error(
			"failed to create zip file",
			zap.String("zip", zipName),
			zap.Error(err),
		)

		return errCreation
	}

	err = outFile.Close()
	if err != nil {
		m.log.Error(
			"failed to close zip file",
			zap.String("zip", zipName),
			zap.Error(err),
		)

		return errFileClosing
	}

	return nil
}

func (m *Manager) GetZipFileData(zipName string) ([]byte, error) {
	path, err := m.buildPathToZipFile(zipName)
	if err != nil {
		return nil, err
	}

	return m.getDataFromFileByPath(path)
}

func (m *Manager) WriteToZipFile(zipName string, file dto.File) error {
	zipPath := filepath.Join(m.zipDir, fmt.Sprintf("%s.zip", zipName))
	outFile, err := os.OpenFile(zipPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		m.log.Error(
			"failed to create zip file",
			zap.String("zip", zipName),
			zap.Error(err),
		)

		return errCreation
	}

	zipWriter := zip.NewWriter(outFile)

	defer func() {
		err = multierr.Combine(err, zipWriter.Close())
		if err != nil {
			m.log.Error(
				"failed to close zip writer",
				zap.String("zip", zipName),
				zap.Error(err),
			)
		}
	}()

	tempPath := filepath.Join(m.fileDir, zipName)
	err = m.writeFileToZip(zipWriter, filepath.Join(tempPath, file.Name), file.Name, file.Ext)
	if err != nil {
		return err
	}

	return nil
}

func (m *Manager) writeFilesToZip(outFile *os.File, zipName string, files []dto.File) (err error) {
	zipWriter := zip.NewWriter(outFile)

	defer func() {
		err = multierr.Combine(err, zipWriter.Close())
		if err != nil {
			m.log.Error(
				"failed to close zip writer",
				zap.String("zip", zipName),
				zap.Error(err),
			)
		}
	}()

	tempPath := filepath.Join(m.fileDir, zipName)
	for _, file := range files {
		err = m.writeFileToZip(zipWriter, filepath.Join(tempPath, file.Name), file.Name, file.Ext)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *Manager) writeFileToZip(zw *zip.Writer, pathToFile, fn, ext string) (err error) {
	writer, err := zw.Create(fn + ext)
	if err != nil {
		m.log.Error(
			"failed to create file inside zip",
			zap.String("fn.ext", fn+ext),
			zap.String("path", pathToFile),
			zap.Error(err),
		)

		return errZipReadWrite
	}

	data, err := m.getDataFromFileByPath(pathToFile)
	if err != nil {
		return err
	}

	var decoded []byte
	err = gob.NewDecoder(bytes.NewBuffer(data)).Decode(&decoded)
	if err != nil {
		m.log.Error(
			"failed to decode file data",
			zap.String("fn.ext", fn+ext),
			zap.String("path", pathToFile),
			zap.Error(err),
		)

		return errGobEncDec
	}

	_, err = writer.Write(decoded)
	if err != nil {
		m.log.Error(
			"failed to write decoded data to file inside zip",
			zap.String("fn.ext", fn+ext),
			zap.String("path", pathToFile),
			zap.Error(err),
		)

		return errZipReadWrite
	}

	return nil
}

func (m *Manager) createTempFileDir() (string, error) {
	tempDir := filepath.Join(m.cfg.AppRootDir, m.cfg.TempFilesDataFolder)
	err := m.createDirByPath(tempDir)
	if err != nil {
		return "", err
	}

	return tempDir, nil
}

func (m *Manager) createTempZipDir() (string, error) {
	zipDir := filepath.Join(m.cfg.AppRootDir, m.cfg.ZipDataFolder)
	err := m.createDirByPath(zipDir)
	if err != nil {
		return "", err
	}

	return zipDir, nil
}

func (m *Manager) createDirByPath(path string) error {
	info, err := os.Stat(path)
	if info != nil && info.IsDir() {
		m.log.Info(
			"temp file dir already exists",
			zap.String("on-path", path),
		)

		return nil
	}

	if err != nil && os.IsNotExist(err) {
		m.log.Info(
			"creating temp file dir",
			zap.String("on path", path),
		)

		err = os.Mkdir(path, os.ModePerm)
		if err != nil {
			m.log.Error(
				"failed to create temp file dir",
				zap.String("on path", path),
				zap.Error(err),
			)

			return errCreation
		}
	}

	return nil
}
