package usecase

import (
	"app/internal/models"
	"app/internal/order"
	"app/pkg/utils"
	"context"
	"time"

	"go.elastic.co/apm"
)

type orderUC struct {
	orderRepo order.Repository
}

func NewOrderUseCase(orderRepo order.Repository) order.UseCase {
	return &orderUC{orderRepo: orderRepo}
}

func (u *orderUC) Create(ctx context.Context, order *models.Order) error {
	span, ctx := apm.StartSpan(ctx, "orderRepo.Create", "custom")
	defer span.End()

	for _, od := range order.OrderList {
		order.TotalPrice += od.PricePerUnit
	}

	now := time.Now().Unix()
	order.CreateAt = now
	order.UpdateAt = now

	return u.orderRepo.Create(ctx, order)
}

func (u *orderUC) GetAll(ctx context.Context, filter *models.OrderFilter, query *utils.PaginationQuery) (*models.OrderList, error) {
	span, ctx := apm.StartSpan(ctx, "orderRepo.GetAll", "custom")
	defer span.End()

	return u.orderRepo.GetAll(ctx, filter, query)
}
