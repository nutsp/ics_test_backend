package http

import (
	"app/config"
	"app/internal/models"
	"app/internal/payment"

	"github.com/gofiber/fiber/v2"
	"go.elastic.co/apm"
)

type paymentHandler struct {
	paymentUC payment.UseCase
	cfg       *config.Config
}

func NewPaymentHandler(paymentUC payment.UseCase, cfg *config.Config) payment.Handlers {
	return &paymentHandler{paymentUC: paymentUC, cfg: cfg}
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

	payment.ReceiptPath = h.cfg.Server.Host + payment.ReceiptPath

	return c.Status(201).JSON(&models.Response{
		Status:  "OK",
		Message: "create payment success",
		Data:    payment,
	})
}

func (h *paymentHandler) UploadReceipt(c *fiber.Ctx) error {
	span, ctx := apm.StartSpan(c.Context(), "orderHandler.UploadReceipt", "custom")
	defer span.End()

	mpf, err := c.MultipartForm()
	if err != nil {
		return c.Status(400).JSON(&models.Response{
			Status:  "FAIL",
			Message: err.Error(),
		})
	}

	path, err := h.paymentUC.UploadReceipt(ctx, mpf)
	if err != nil {
		return c.Status(400).JSON(&models.Response{
			Status:  "FAIL",
			Message: err.Error(),
		})
	}

	return c.Status(200).JSON(&models.Response{
		Status: "OK",
		Data: fiber.Map{
			"path": path,
		},
	})
}
