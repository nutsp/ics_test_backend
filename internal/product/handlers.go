package product

import "github.com/gofiber/fiber/v2"

type Handlers interface {
	Create(c *fiber.Ctx) error
	GetAll(c *fiber.Ctx) error
}
