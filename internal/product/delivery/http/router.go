package http

import (
	"app/internal/product"

	"github.com/gofiber/fiber/v2"
)

func MapProductRoute(r fiber.Router, h product.Handlers) {
	r.Post("/", h.Create)
	r.Get("/", h.GetAll)
}
