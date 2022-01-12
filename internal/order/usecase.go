package order

import (
	"app/internal/models"
	"app/pkg/utils"
	"context"
)

type UseCase interface {
	Create(ctx context.Context, order *models.Order) error
	GetAll(ctx context.Context, filter *models.OrderFilter, query *utils.PaginationQuery) (*models.OrderList, error)
	GetByUserID(ctx context.Context, filter *models.OrderFilter, query *utils.PaginationQuery) (*models.OrderList, error)
	GetByID(ctx context.Context, id int) (*models.OrderBase, error)
}
