package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
)

// AuthReq middleware
func AdminAuthReq() func(c *fiber.Ctx) error {
	cfg := basicauth.Config{
		Users: map[string]string{
			"admin": "123456",
		},
		Realm: "Forbidden",
		Authorizer: func(u, p string) bool {
			if u == "admin" && p == "123456" {
				return true
			}
			return false
		},
		ContextUsername: "_user",
		ContextPassword: "_pass",
	}

	return basicauth.New(cfg)
}

func UserAuthReq() func(c *fiber.Ctx) error {
	cfg := basicauth.Config{
		Users: map[string]string{
			"user":  "123456",
			"admin": "123456",
		},
	}

	return basicauth.New(cfg)
}
