package user

import (
	"app/internal/core/ports"
	"app/internal/core/services/user"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ohmspeed777/go-pkg/corex"
)

type Handler struct {
	user user.IUserService
}

func NewHandler(user user.IUserService) ports.IUserHandler {
	return &Handler{
		user: user,
	}
}

func (h *Handler) GetMe(c echo.Context) error {
	res, err := h.user.GetMe(corex.NewFromEchoContext(c))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, res)
}

func (h *Handler) GetMyOrderHistory(c echo.Context) error {
	res, err := h.user.GetMyOrderHistory(corex.NewFromEchoContext(c))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, res)
}

