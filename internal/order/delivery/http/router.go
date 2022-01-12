package http

import (
	"app/internal/order"
	"app/pkg/middleware"

	"github.com/gofiber/fiber/v2"
)

func MapOrderRoute(r fiber.Router, h order.Handlers) {
	r.Use(middleware.UserAuthReq())
	r.Post("/", h.Create)
	r.Get("/", h.GetAll)

	r.Use(middleware.AdminAuthReq())
	r.Get("/list", h.GetAll)
}
