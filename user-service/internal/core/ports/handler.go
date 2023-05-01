package ports

import "github.com/labstack/echo/v4"

type IUserHandler interface {
	GetMe(c echo.Context) error
}
