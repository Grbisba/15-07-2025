package manager

import (
	"github.com/Grbisba/15-07-2025/internal/model/dto"
)

type Manager interface {
	CreateFile(fileData []byte, dirName string, fn string) error
	CreateZipFile(zipName string) error
	GetZipFileData(zipName string) ([]byte, error)
	WriteToZipFile(zipName string, file dto.File) error
}
