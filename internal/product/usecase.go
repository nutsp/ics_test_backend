package product

import (
	"app/internal/models"
	"app/pkg/utils"
	"context"
)

type UseCase interface {
	Create(ctx context.Context, product *models.Product) error
	GetAll(ctx context.Context, filter *models.ProductFilter, query *utils.PaginationQuery) (*models.ProductList, error)
}
