package fiber

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"

	"github.com/Grbisba/15-07-2025/internal/pkg/ownerr"
)

func (ctrl *Controller) handleErr(err error) error {
	switch {
	case err == nil:
		return fiber.NewError(fiber.StatusOK, "OK")
	case errors.Is(err, ownerr.ErrUnknown):
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	case errors.Is(err, ownerr.ErrBadRequest):
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	case errors.Is(err, ownerr.ErrNotFound):
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	case errors.Is(err, ownerr.ErrAlreadyExist):
		return fiber.NewError(fiber.StatusConflict, err.Error())
	case errors.Is(err, ownerr.ErrUnauthorized):
		return fiber.NewError(fiber.StatusUnauthorized, err.Error())
	case errors.Is(err, ownerr.ErrForbidden):
		return fiber.NewError(fiber.StatusForbidden, err.Error())
	case errors.Is(err, ownerr.ErrServiceUnavailable):
		return fiber.NewError(fiber.StatusServiceUnavailable, err.Error())
	case errors.Is(err, ownerr.ErrMovedPermanently):
		return fiber.NewError(fiber.StatusMovedPermanently, err.Error())
	case errors.Is(err, ownerr.ErrInternal):
		fallthrough
	default:
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
}
