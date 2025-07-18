package fiber

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/Grbisba/15-07-2025/internal/controller/http/model"
	"github.com/Grbisba/15-07-2025/internal/model/dto"
)

func (ctrl *Controller) Health(c *fiber.Ctx) error {
	return c.SendString("ok")
}

func (ctrl *Controller) GetStatus(c *fiber.Ctx) error {
	ctrl.log.Info("handle get task status")

	ctx := c.Context()

	authed, err := ctrl.service.CheckUserAuth(ctx, c.Get(fiber.HeaderAuthorization))
	if err != nil {
		return ctrl.handleErr(err)
	}

	if authed == nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	rawID := c.Params("task_id")
	id, err := uuid.Parse(rawID)
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	resp, err := ctrl.service.GetStatus(ctx, id, authed.ID)
	if err != nil {
		return ctrl.handleErr(err)
	}

	return c.Status(fiber.StatusOK).JSON(model.GetStatusResponse{
		UserID:        resp.UserID,
		Status:        resp.Status,
		URL:           resp.URL,
		UploadedFiles: resp.UploadedFiles,
	})
}

func (ctrl *Controller) CreateTask(c *fiber.Ctx) error {
	ctrl.log.Info("handle create task")

	authed, err := ctrl.service.CheckUserAuth(c.Context(), c.Get(fiber.HeaderAuthorization))
	if err != nil {
		return ctrl.handleErr(err)
	}

	if authed == nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	task, err := ctrl.service.CreateTask(c.Context(), authed.ID)
	if err != nil {
		return ctrl.handleErr(err)
	}

	return c.Status(fiber.StatusCreated).JSON(model.TaskResponse{
		ID:     task.ID,
		Status: task.Status,
	})
}

func (ctrl *Controller) UploadFile(c *fiber.Ctx) error {
	ctrl.log.Info("handle upload file")

	authed, err := ctrl.service.CheckUserAuth(c.Context(), c.Get(fiber.HeaderAuthorization))
	if err != nil {
		return ctrl.handleErr(err)
	}

	if authed == nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	rawID := c.Params("task_id")
	id, err := uuid.Parse(rawID)
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	req := new(model.UploadFileRequest)
	if err = c.BodyParser(req); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	task, err := ctrl.service.UploadFile(c.Context(), dto.UploadFile{
		UserID: authed.ID,
		TaskID: id,
		URLs:   req.URLs,
	})
	if err != nil {
		return ctrl.handleErr(err)
	}

	return c.Status(fiber.StatusCreated).JSON(model.TaskResponse{
		ID:     task.ID,
		Status: task.Status,
	})
}

func (ctrl *Controller) Download(c *fiber.Ctx) error {
	ctrl.log.Info("handle download file")

	rawID := c.Params("task_id")
	id, err := uuid.Parse(rawID)
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	resp, err := ctrl.service.Download(c.Context(), id)
	if err != nil {
		return ctrl.handleErr(err)
	}

	c.Set(fiber.HeaderContentType, "application/zip")
	c.Set(fiber.HeaderContentDisposition, fmt.Sprintf("attachment; filename=\"%s.zip\"", resp.Filename))
	return c.Send(resp.RawZip)
}

func (ctrl *Controller) CreateUser(c *fiber.Ctx) error {
	ctrl.log.Info("handle create user")

	req := new(model.CreateUserRequest)
	if err := c.BodyParser(req); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	resp, err := ctrl.service.CreateUser(c.Context(), dto.CreateUser{
		Login:    req.Login,
		Password: req.Password,
	})
	if err != nil {
		return ctrl.handleErr(err)
	}

	ctrl.log.Debug("successfully handle create consumer")

	return c.Status(fiber.StatusCreated).JSON(model.CreateUserResponse{
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
	})
}

func (ctrl *Controller) Login(c *fiber.Ctx) error {
	ctrl.log.Info("handle login user")

	req := new(model.CreateUserRequest)
	if err := c.BodyParser(req); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	resp, err := ctrl.service.Login(c.Context(), dto.Login{
		Login:    req.Login,
		Password: req.Password,
	})
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(model.CreateUserResponse{
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
	})
}

func (ctrl *Controller) RefreshToken(c *fiber.Ctx) error {
	resp, err := ctrl.service.RefreshToken(c.Context(), c.Get(fiber.HeaderAuthorization))
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(model.RefreshTokenResponse{
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
	})
}
