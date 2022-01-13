package http

import (
	"app/internal/payment"

	"github.com/gofiber/fiber/v2"
)

func MapPaymentRoute(r fiber.Router, h payment.Handlers) {
	r.Post("/", h.Create)
	r.Post("/upload", h.UploadReceipt)
}
