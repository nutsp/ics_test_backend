package http

import (
	"app/internal/models"
	"app/internal/payment"

	"github.com/gofiber/fiber/v2"
	"go.elastic.co/apm"
)

type paymentHandler struct {
	paymentUC payment.UseCase
}

func NewPaymentHandler(paymentUC payment.UseCase) payment.Handlers {
	return &paymentHandler{paymentUC: paymentUC}
}

func (h *paymentHandler) Create(c *fiber.Ctx) error {
	span, ctx := apm.StartSpan(c.Context(), "orderHandler.Create", "custom")
	defer span.End()

	payment := new(models.Payment)
	if err := c.BodyParser(payment); err != nil {
		return c.Status(400).JSON(&models.Response{
			Status:  "FAIL",
			Message: "body parser payment fail",
			Error:   err.Error(),
		})
	}

	if err := h.paymentUC.Create(ctx, payment); err != nil {
		return c.Status(400).JSON(&models.Response{
			Status:  "FAIL",
			Message: "create payment fail",
			Error:   err.Error(),
		})
	}

	return c.Status(201).JSON(&models.Response{
		Status:  "OK",
		Message: "create payment success",
		Data:    payment,
	})
}
