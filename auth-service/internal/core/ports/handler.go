package ports

import "github.com/labstack/echo/v4"

type IAuthHandler interface {
	SignUp(c echo.Context) error
	SignIn(c echo.Context) error
}
