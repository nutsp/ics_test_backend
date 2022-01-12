package http

import (
	"app/internal/order"

	"github.com/gofiber/fiber/v2"
)

func MapOrderRoute(r fiber.Router, h order.Handlers) {
	r.Post("/", h.Create)
	r.Get("/", h.GetAll)
}
