package http

import (
	"app/internal/order"
	"app/pkg/middleware"

	"github.com/gofiber/fiber/v2"
)

func MapOrderRoute(r fiber.Router, h order.Handlers) {
	r.Use(middleware.UserAuthReq())
	r.Post("/", h.Create)
	r.Get("/user/all", h.GetByUserID)
	r.Get("/:id", h.GetByID)
	r.Put("/confirm/:id", h.UpdateConfirm)
	r.Put("/cancel/:id", h.UpdateCancel)

	r.Use(middleware.AdminAuthReq())
	r.Get("/", h.GetAll)
}
