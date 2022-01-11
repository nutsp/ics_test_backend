package server

import (
	productHttp "app/internal/product/delivery/http"
	productRepository "app/internal/product/repository"
	productUseCase "app/internal/product/usecase"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func (s *Server) MapHandlers(app *fiber.App) error {

	config := cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
		AllowCredentials: true,
	}
	app.Use(cors.New(config))
	app.Static("/storage", "./storage", fiber.Static{
		Compress:      true,
		ByteRange:     true,
		Browse:        true,
		CacheDuration: 10 * time.Second,
		MaxAge:        3600,
	})

	v1 := app.Group("/v1")

	// Products
	productGroup := v1.Group("/product")
	productRepo := productRepository.NewProductRepository(s.db)
	productUC := productUseCase.NewProductUseCase(productRepo)
	productHandler := productHttp.NewProductHandler(productUC)
	productHttp.MapProductRoute(productGroup, productHandler)

	// Orders

	// Payments

	return nil
}