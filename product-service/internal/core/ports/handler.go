package ports

import "github.com/labstack/echo/v4"

type IProductHandler interface {
	GetAll(c echo.Context) error
	GetOne(c echo.Context) error
	Create(c echo.Context) error
}
