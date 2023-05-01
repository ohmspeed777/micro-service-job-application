package ports

import "github.com/labstack/echo/v4"

type IOrderHandler interface {
	Create(c echo.Context) error
	Cancel(c echo.Context) error
	GetOne(c echo.Context) error
}
