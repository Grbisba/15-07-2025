package production

import (
	"io"
	"net/http"

	"github.com/pkg/errors"

	"github.com/Grbisba/15-07-2025/internal/model/dto"
	"github.com/Grbisba/15-07-2025/internal/pkg/ownerr"
	"github.com/Grbisba/15-07-2025/internal/pkg/status"
)

func (s *Service) fetchDataFromURL(url string, file *dto.File) ([]byte, error) {
	file.URL = url

	resp, err := http.Get(url)
	if err != nil {
		return nil, ownerr.NewError(
			ownerr.ErrInternal,
			errors.Wrap(err, "failed to download file"),
		)
	}

	fileData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, ownerr.NewError(
			ownerr.ErrInternal,
			errors.Wrap(err, "failed to download file"),
		)
	}

	switch http.DetectContentType(fileData) {
	case "application/pdf":
		file.Ext = ".pdf"
		file.Status = status.Uploaded
	case "image/jpeg":
		file.Ext = ".jpeg"
		file.Status = status.Uploaded
	default:
		file.Ext = "none"
		file.Status = status.Unsupported
	}

	return fileData, nil
}
