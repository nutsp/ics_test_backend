package http

import (
	"app/internal/models"
	"app/internal/product"
	"app/pkg/utils"
	"log"

	"github.com/gofiber/fiber/v2"
	"go.elastic.co/apm"
)

type productHandler struct {
	productUC product.UseCase
}

func NewProductHandler(productUC product.UseCase) product.Handlers {
	return &productHandler{productUC: productUC}
}

func (h *productHandler) Create(c *fiber.Ctx) error {
	span, ctx := apm.StartSpan(c.Context(), "productHandler.Create", "custom")
	defer span.End()

	userID := c.Get("Authorization")
	log.Println(userID)
	product := new(models.Product)
	if err := c.BodyParser(product); err != nil {
		return c.Status(400).JSON(&models.Response{
			Status:  "FAIL",
			Message: "body parser product fail",
			Error:   err.Error(),
		})
	}

	if err := h.productUC.Create(ctx, product); err != nil {
		return c.Status(400).JSON(&models.Response{
			Status:  "FAIL",
			Message: "create product fail",
			Error:   err.Error(),
		})
	}

	return c.Status(201).JSON(&models.Response{
		Status:  "OK",
		Message: "create product success",
		Data:    product,
	})
}

func (h *productHandler) GetAll(c *fiber.Ctx) error {
	span, ctx := apm.StartSpan(c.Context(), "productHandler.GetAll", "custom")
	defer span.End()

	userID := c.Get("Authorization")
	log.Println(userID)

	ft := new(models.ProductFilter)
	pq := new(utils.PaginationQuery)
	_ = c.QueryParser(ft)
	_ = c.QueryParser(pq)

	productList, err := h.productUC.GetAll(ctx, ft, pq)
	if err != nil {
		return c.Status(400).JSON(&models.Response{
			Status:  "FAIL",
			Message: "list product fail",
			Error:   err.Error(),
		})
	}

	return c.Status(200).JSON(&models.Response{
		Status:  "OK",
		Message: "list product success",
		Data:    productList,
	})
}
