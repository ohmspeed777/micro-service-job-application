package order

import (
	"app/internal/core/domain"
	"app/internal/core/ports"
	"app/internal/core/services/order"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ohmspeed777/go-pkg/corex"
	"github.com/ohmspeed777/go-pkg/errorx"
)

type Handler struct {
	order order.IOrderService
}

func NewHandler(order order.IOrderService) ports.IOrderHandler {
	return &Handler{
		order: order,
	}
}

func (h *Handler) Create(c echo.Context) error {
	req := order.CreateReq{}
	if err := c.Bind(&req); err != nil {
		return errorx.NewInvalidRequest(err)
	}

	res, err := h.order.Create(corex.NewFromEchoContext(c), req)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, res)
}

func (h *Handler) Cancel(c echo.Context) error {
	req := domain.GetOneReq{}
	if err := c.Bind(&req); err != nil {
		return errorx.NewInvalidRequest(err)
	}

	res, err := h.order.Cancel(corex.NewFromEchoContext(c), req)
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

	res, err := h.order.GetOne(corex.NewFromEchoContext(c), req)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, res)
}


