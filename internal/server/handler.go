package server

import (
	productHttp "app/internal/product/delivery/http"
	productRepository "app/internal/product/repository"
	productUseCase "app/internal/product/usecase"

	orderHttp "app/internal/order/delivery/http"
	orderRepository "app/internal/order/repository"
	orderUseCase "app/internal/order/usecase"

	paymentHttp "app/internal/payment/delivery/http"
	paymentRepository "app/internal/payment/repository"
	paymentUseCase "app/internal/payment/usecase"
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
	orderGroup := v1.Group("/order")
	orderRepo := orderRepository.NewOrderRepository(s.db)
	orderUC := orderUseCase.NewOrderUseCase(orderRepo)
	orderHandler := orderHttp.NewOrderHandler(orderUC)
	orderHttp.MapOrderRoute(orderGroup, orderHandler)

	// Payments
	paymentGroup := v1.Group("/payment")
	paymentRepo := paymentRepository.NewPaymentRepository(s.db)
	paymentUC := paymentUseCase.NewPaymentUseCase(paymentRepo)
	paymentHandler := paymentHttp.NewPaymentHandler(paymentUC, s.cfg)
	paymentHttp.MapPaymentRoute(paymentGroup, paymentHandler)
	return nil
}
