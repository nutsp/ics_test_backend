package http

import (
	"app/internal/models"
	"app/internal/order"
	"app/pkg/middleware"
	"app/pkg/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"go.elastic.co/apm"
)

type orderHandler struct {
	orderUC order.UseCase
}

func NewOrderHandler(orderUC order.UseCase) order.Handlers {
	return &orderHandler{orderUC: orderUC}
}

func (h *orderHandler) Create(c *fiber.Ctx) error {
	span, ctx := apm.StartSpan(c.Context(), "orderHandler.Create", "custom")
	defer span.End()

	order := new(models.Order)
	if err := c.BodyParser(order); err != nil {
		return c.Status(400).JSON(&models.Response{
			Status:  "FAIL",
			Message: "body parser order fail",
			Error:   err.Error(),
		})
	}

	order.UserID = middleware.GetUserID(c.Get("Authorization"))

	if err := h.orderUC.Create(ctx, order); err != nil {
		return c.Status(400).JSON(&models.Response{
			Status:  "FAIL",
			Message: "create order fail",
			Error:   err.Error(),
		})
	}

	return c.Status(201).JSON(&models.Response{
		Status:  "OK",
		Message: "create order success",
		Data:    order,
	})
}

func (h *orderHandler) GetAll(c *fiber.Ctx) error {
	span, ctx := apm.StartSpan(c.Context(), "orderHandler.GetAll", "custom")
	defer span.End()

	ft := new(models.OrderFilter)
	pq := new(utils.PaginationQuery)
	_ = c.QueryParser(ft)
	_ = c.QueryParser(pq)

	orderList, err := h.orderUC.GetAll(ctx, ft, pq)
	if err != nil {
		return c.Status(400).JSON(&models.Response{
			Status:  "FAIL",
			Message: "list order fail",
			Error:   err.Error(),
		})
	}

	return c.Status(200).JSON(&models.Response{
		Status:  "OK",
		Message: "list order success",
		Data:    orderList,
	})
}

func (h *orderHandler) GetByUserID(c *fiber.Ctx) error {
	span, ctx := apm.StartSpan(c.Context(), "orderHandler.GetByUserID", "custom")
	defer span.End()

	ft := new(models.OrderFilter)
	pq := new(utils.PaginationQuery)
	_ = c.QueryParser(ft)
	_ = c.QueryParser(pq)

	ft.UserID = middleware.GetUserID(c.Get("Authorization"))
	orderList, err := h.orderUC.GetByUserID(ctx, ft, pq)
	if err != nil {
		return c.Status(400).JSON(&models.Response{
			Status:  "FAIL",
			Message: "list order fail",
			Error:   err.Error(),
		})
	}

	return c.Status(200).JSON(&models.Response{
		Status:  "OK",
		Message: "list order success",
		Data:    orderList,
	})
}

func (h *orderHandler) GetByID(c *fiber.Ctx) error {
	span, ctx := apm.StartSpan(c.Context(), "orderHandler.GetByID", "custom")
	defer span.End()

	id, _ := strconv.Atoi(c.Params("id"))
	// ft.UserID = middleware.GetUserID(c.Get("Authorization"))
	orderList, err := h.orderUC.GetByID(ctx, id)
	if err != nil {
		return c.Status(400).JSON(&models.Response{
			Status:  "FAIL",
			Message: "list order fail",
			Error:   err.Error(),
		})
	}

	return c.Status(200).JSON(&models.Response{
		Status:  "OK",
		Message: "list order success",
		Data:    orderList,
	})
}
