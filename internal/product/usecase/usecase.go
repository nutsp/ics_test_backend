package usecase

import (
	"app/internal/models"
	"app/internal/product"
	"app/pkg/utils"
	"context"
	"time"

	"go.elastic.co/apm"
)

type productUC struct {
	productRepo product.Repository
}

func NewProductUseCase(productRepo product.Repository) product.UseCase {
	return &productUC{productRepo: productRepo}
}

func (u *productUC) Create(ctx context.Context, product *models.Product) error {
	span, ctx := apm.StartSpan(ctx, "productUC.Create", "custom")
	defer span.End()

	now := time.Now().Unix()
	product.CreateAt = now
	product.UpdateAt = now

	return u.productRepo.Create(ctx, product)
}

func (u *productUC) GetAll(ctx context.Context, filter *models.ProductFilter, query *utils.PaginationQuery) (*models.ProductList, error) {
	span, ctx := apm.StartSpan(ctx, "productUC.GetAll", "custom")
	defer span.End()

	return u.productRepo.GetAll(ctx, filter, query)
}
