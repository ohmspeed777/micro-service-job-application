package auth

import (
	"app/internal/core/ports"
	"app/internal/core/services/auth"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ohmspeed777/go-pkg/corex"
	"github.com/ohmspeed777/go-pkg/errorx"
)

type Handler struct {
	auth auth.IAuthService
}

func NewHandler(auth auth.IAuthService) ports.IAuthHandler {
	return &Handler{
		auth: auth,
	}
}

func (h *Handler) SignUp(c echo.Context) error {
	req := auth.SignUpReq{}
	if err := c.Bind(&req); err != nil {
		return errorx.New(http.StatusBadRequest, "can not bind request", err)
	}

	res, err := h.auth.SignUp(corex.NewFromEchoContext(c), req)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, res)
}


func (h *Handler) SignIn(c echo.Context) error {
	req := auth.SignInReq{}
	if err := c.Bind(&req); err != nil {
		return errorx.New(http.StatusBadRequest, "can not bind request", err)
	}

	res, err := h.auth.SignIn(corex.NewFromEchoContext(c), req)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, res)
}