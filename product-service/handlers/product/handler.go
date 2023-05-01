package product

import (
	"app/internal/core/domain"
	"app/internal/core/ports"
	"app/internal/core/services/product"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ohmspeed777/go-pkg/corex"
	"github.com/ohmspeed777/go-pkg/errorx"
)

type Handler struct {
	product product.IProductService
}

func NewHandler(product product.IProductService) ports.IProductHandler {
	return &Handler{
		product: product,
	}
}

func (h *Handler) Create(c echo.Context) error {
	req := product.CreateReq{}
	if err := c.Bind(&req); err != nil {
		return errorx.NewInvalidRequest(err)
	}

	res, err := h.product.Create(corex.NewFromEchoContext(c), req)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, res)
}

func (h *Handler) GetOne(c echo.Context) error {
	req := domain.GetOneReq{}
	if err := c.Bind(&req); err != nil {
		return errorx.NewInvalidRequest(err)
	}

	res, err := h.product.GetOne(corex.NewFromEchoContext(c), req)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, res)
}


func (h *Handler) GetAll(c echo.Context) error {
	req := product.GetAllReq{}
	if err := c.Bind(&req); err != nil {
		return errorx.NewInvalidRequest(err)
	}

	res, err := h.product.GetAll(corex.NewFromEchoContext(c), req)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, res)
}