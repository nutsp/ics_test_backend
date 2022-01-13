package payment

import "github.com/gofiber/fiber/v2"

type Handlers interface {
	Create(c *fiber.Ctx) error
	UploadReceipt(c *fiber.Ctx) error
}
