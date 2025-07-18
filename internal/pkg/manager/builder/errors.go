package builder

import (
	"github.com/pkg/errors"
)

var (
	errCreation     = errors.New("error while creating directory or file")
	errOsStat       = errors.New("error while getting directory or file info")
	errNotExist     = errors.New("file or directory does not exist")
	errGobEncDec    = errors.New("error while decoding or encoding binary file")
	errNilReference = errors.New("provided nil dependency")
	errFileReading  = errors.New("error while reading file")
	errFileClosing  = errors.New("error while closing file")
	errZipReadWrite = errors.New("error while reading or writing zip file")
	errBuildPath    = errors.New("error while building path")
)
