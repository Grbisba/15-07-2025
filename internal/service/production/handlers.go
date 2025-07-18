package production

import (
	"context"
	"path"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/Grbisba/15-07-2025/internal/model/dto"
	"github.com/Grbisba/15-07-2025/internal/pkg/ownerr"
	"github.com/Grbisba/15-07-2025/internal/pkg/status"
)

// GetStatus -> status
func (s *Service) GetStatus(_ context.Context, id uuid.UUID, userID uuid.UUID) (*dto.GetStatus, error) {
	user, ok := s.store.Get(userID.String())
	if !ok {
		return nil, ownerr.NewError(
			ownerr.ErrNotFound,
			errors.New("user not found or not exists"),
		)
	}

	task, ok := s.tasks.Get(id)
	if !ok {
		return nil, ownerr.NewError(
			ownerr.ErrNotFound,
			errors.New("task not found or not exists"),
		)
	}

	resp := &dto.GetStatus{
		UserID:        user.ID,
		Status:        task.Status,
		UploadedFiles: len(task.Files),
	}

	if resp.UploadedFiles >= 3 {
		resp.URL = s.createURL(task.ID)
	}

	return resp, nil
}

// CreateTask -> status with task id
func (s *Service) CreateTask(_ context.Context, userID uuid.UUID) (*dto.Task, error) {
	user, ok := s.store.Get(userID.String())
	if !ok {
		return nil, ownerr.NewError(
			ownerr.ErrNotFound,
			errors.New("user not found or not exists"),
		)
	}

	if user.TotalTasks >= s.cfg.Loader.MaxConcurrentTasks {
		return nil, ownerr.NewError(
			ownerr.ErrServiceUnavailable,
			errors.New("max concurrent tasks exceeded"),
		)
	}

	task := dto.Task{
		UserID: user.ID,
		ID:     uuid.New(),
		Status: status.Created,
		Files:  make([]dto.File, 0, 3),
	}

	err := s.dm.CreateZipFile(task.ID.String())
	if err != nil {
		return nil, ownerr.NewError(
			ownerr.ErrInternal,
			errors.Wrap(err, "failed to create zip file"),
		)
	}

	s.tasks.Set(task.ID, task)

	return &task, nil
}

// UploadFile -> status
func (s *Service) UploadFile(_ context.Context, req dto.UploadFile) (*dto.Task, error) {
	task, ok := s.tasks.Get(req.TaskID)
	if !ok {
		return nil, ownerr.NewError(
			ownerr.ErrNotFound,
			errors.New("task not found or not exists"),
		)
	}

	if len(task.Files) >= s.cfg.Loader.MaxFilesPerTask {
		return nil, ownerr.NewError(
			ownerr.ErrForbidden,
			errors.New("max files per task exceeded"),
		)
	}

	for _, url := range req.URLs {
		if len(task.Files) >= s.cfg.Loader.MaxFilesPerTask {
			task.Status = status.Completed
			s.tasks.Set(task.ID, task)
			return &task, nil
		}

		file := dto.File{
			Name: uuid.New().String(),
			URL:  url,
		}

		fileData, err := s.fetchDataFromURL(url, &file)
		if err != nil {
			s.log.Warn(
				"failed to fetch data from url",
				zap.Error(err),
			)

			file.Ext = "none"
			file.Status = status.BadURL
			task.Files = append(task.Files, file)
			continue
		}

		err = s.dm.CreateFile(fileData, req.TaskID.String(), file.Name)
		if err != nil {
			file.Status = status.Failed
			task.Files = append(task.Files, file)
			continue
		}

		err = s.dm.WriteToZipFile(task.ID.String(), file)
		if err != nil {
			file.Status = status.Failed
			task.Files = append(task.Files, file)
			continue
		}

		task.Files = append(task.Files, file)
	}

	if len(task.Files) > 0 {
		task.Status = status.InProcess
	}

	s.tasks.Set(task.ID, task)

	return &task, nil
}

// Download -> status
func (s *Service) Download(_ context.Context, id uuid.UUID) (*dto.Download, error) {
	var err error
	download := new(dto.Download)

	task, ok := s.tasks.Get(id)
	if !ok {
		s.log.Error(
			"task not found",
			zap.String("id", id.String()),
		)

		return nil, errors.New("task not found")
	}

	download.RawZip, err = s.dm.GetZipFileData(id.String())
	if err != nil {
		s.log.Error(
			"download get zip file error",
			zap.Error(err),
		)

		return nil, err
	}

	task.Status = status.Downloaded
	download.Status = task.Status
	download.Filename = id.String()

	s.tasks.Set(task.ID, task)

	return download, nil
}

func (s *Service) CreateUser(_ context.Context, req dto.CreateUser) (*dto.Token, error) {
	_, ok := s.store.Get(req.Login)
	if ok {
		s.log.Warn(
			"user already exists",
			zap.String("login", req.Login),
		)

		return nil, ownerr.NewError(
			ownerr.ErrAlreadyExist,
			errors.New("user already exists"),
		)
	}

	password, err := s.enc.EncryptPassword(req.Password)
	if err != nil {
		s.log.Error(
			"failed to encrypt password",
			zap.Error(err),
		)

		return nil, ownerr.NewError(
			ownerr.ErrInternal,
			errors.Wrap(err, "failed to encrypt password"),
		)
	}

	user := dto.User{
		ID:          uuid.New(),
		Login:       req.Login,
		EncPassword: password,
		TotalTasks:  0,
	}

	s.store.Set(user.Login, user.ID, user)

	access, refresh, err := s.prv.CreateAccessAndRefreshTokenForUser(user.ID)
	if err != nil {
		s.log.Error(
			"failed to create a couples of tokens",
			zap.Error(err),
		)

		return nil, ownerr.NewError(
			ownerr.ErrInternal,
			errors.Wrap(err, "failed to create a couples of tokens"),
		)
	}

	return &dto.Token{
		AccessToken:  access,
		RefreshToken: refresh,
	}, nil
}

func (s *Service) Login(_ context.Context, req dto.Login) (*dto.Token, error) {
	user, ok := s.store.Get(req.Login)
	if !ok {
		s.log.Error(
			"failed to get consumer by login: not found",
			zap.String("login", req.Login),
		)

		return nil, ownerr.NewError(
			ownerr.ErrNotFound,
			errors.New("failed to get consumer by login"),
		)
	}

	encryptedPassword, err := s.enc.EncryptPassword(req.Password)
	if err != nil {
		s.log.Error(
			"failed to encrypt password",
			zap.Error(err),
		)

		return nil, ownerr.NewError(
			ownerr.ErrInternal,
			errors.Wrap(err, "failed to encrypt password"),
		)
	}

	err = s.enc.CompareHashAndPassword(user.EncPassword, encryptedPassword)
	if err != nil {
		s.log.Error(
			"failed to compare password",
			zap.Error(err),
		)

		return nil, ownerr.NewError(
			ownerr.ErrForbidden,
			errors.Wrap(err, "failed to compare old passwords"),
		)
	}

	access, refresh, err := s.prv.CreateAccessAndRefreshTokenForUser(user.ID)
	if err != nil {
		s.log.Error(
			"failed to create a couples of tokens",
			zap.Error(err),
		)

		return nil, ownerr.NewError(
			ownerr.ErrInternal,
			errors.Wrap(err, "failed to create a couples of tokens"),
		)
	}

	return &dto.Token{
		AccessToken:  access,
		RefreshToken: refresh,
	}, nil
}

func (s *Service) RefreshToken(_ context.Context, token string) (*dto.Token, error) {
	data, err := s.prv.GetDataFromToken(token)
	if err != nil {
		s.log.Error(
			"failed to get consumer data from token",
			zap.Error(err),
		)

		return nil, ownerr.NewError(
			ownerr.ErrUnauthorized,
			errors.Wrap(err, "failed to get consumer data from token"),
		)
	}

	if data.IsAccess {
		s.log.Error(
			"not acceptable token",
			zap.Bool("is_access", data.IsAccess),
		)

		return nil, ownerr.NewError(
			ownerr.ErrForbidden,
			errors.New("not acceptable token"),
		)
	}

	access, refresh, err := s.prv.CreateAccessAndRefreshTokenForUser(data.ID)
	if err != nil {
		s.log.Error(
			"failed to create a couples of tokens",
			zap.Error(err),
		)

		return nil, ownerr.NewError(
			ownerr.ErrInternal,
			errors.Wrap(err, "failed to create a couples of tokens"),
		)
	}

	return &dto.Token{
		AccessToken:  access,
		RefreshToken: refresh,
	}, nil
}

func (s *Service) CheckUserAuth(_ context.Context, token string) (*dto.UserDataInToken, error) {
	data, err := s.prv.GetDataFromToken(token)
	if err != nil {
		s.log.Error(
			"failed to get consumer data from token",
			zap.Error(err),
		)

		return nil, ownerr.NewError(
			ownerr.ErrUnauthorized,
			errors.Wrap(err, "failed to get consumer data from token"),
		)
	}

	switch {
	case data == nil:
		s.log.Error("failed to get consumer data from token: collected nil data")

		return nil, ownerr.NewError(
			ownerr.ErrForbidden,
			errors.New("failed to get consumer data from token: collected nil data"),
		)
	case !data.IsAccess:
		s.log.Error("failed to get consumer data from token: should be access")

		return nil, ownerr.NewError(
			ownerr.ErrForbidden,
			errors.New("failed to get consumer data from token: should be access"),
		)
	default:
		return data, nil
	}
}

func (s *Service) createURL(id uuid.UUID) string {
	var b strings.Builder
	b.Grow(128)
	b.WriteString("http://")
	b.WriteString(path.Join(
		s.cfg.Controller.Host+":"+strconv.Itoa(s.cfg.Controller.Port),
		"tasks", "download", id.String(),
	))
	return b.String()
}
