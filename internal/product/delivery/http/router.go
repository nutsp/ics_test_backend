package http

import (
	"app/internal/product"
	"app/pkg/middleware"

	"github.com/gofiber/fiber/v2"
)

func MapProductRoute(r fiber.Router, h product.Handlers) {

	// Private
	{
		r.Use(middleware.AdminAuthReq())
		r.Post("/", h.Create)
		r.Get("/", h.GetAll)
	}
}
