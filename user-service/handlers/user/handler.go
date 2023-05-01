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

// func (h *Handler) SignIn(c echo.Context) error {
// 	req := auth.SignInReq{}
// 	if err := c.Bind(&req); err != nil {
// 		return errorx.NewInvalidRequest(err)
// 	}

// 	res, err := h.auth.SignIn(corex.NewFromEchoContext(c), req)
// 	if err != nil {
// 		return err
// 	}

// 	return c.JSON(http.StatusCreated, res)
// }
