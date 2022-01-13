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
	span, ctx := apm.StartSpan(ctx, "orderUC.Create", "custom")
	defer span.End()

	for _, od := range order.OrderList {
		order.TotalPrice += od.PricePerUnit * float64(od.Amount)
	}

	order.StatusID = 1 // Pending Payment

	now := time.Now().Unix()
	order.CreateAt = now
	order.UpdateAt = now

	return u.orderRepo.Create(ctx, order)
}

func (u *orderUC) GetAll(ctx context.Context, filter *models.OrderFilter, query *utils.PaginationQuery) (*models.OrderList, error) {
	span, ctx := apm.StartSpan(ctx, "orderUC.GetAll", "custom")
	defer span.End()

	return u.orderRepo.GetAll(ctx, filter, query)
}

func (u *orderUC) GetByUserID(ctx context.Context, filter *models.OrderFilter, query *utils.PaginationQuery) (*models.OrderList, error) {
	span, ctx := apm.StartSpan(ctx, "orderUC.GetByUserID", "custom")
	defer span.End()

	return u.orderRepo.GetByUserID(ctx, filter, query)
}

func (u *orderUC) GetByID(ctx context.Context, id int) (*models.OrderBase, error) {
	span, ctx := apm.StartSpan(ctx, "orderUC.GetByID", "custom")
	defer span.End()

	return u.orderRepo.GetByID(ctx, id)
}

func (u *orderUC) UpdateCancel(ctx context.Context, id int) error {
	span, ctx := apm.StartSpan(ctx, "orderUC.UpdateCancel", "custom")
	defer span.End()

	statusID := 4

	return u.orderRepo.UpdateCancel(ctx, id, statusID)
}

func (u *orderUC) UpdateConfirm(ctx context.Context, id int) error {
	span, ctx := apm.StartSpan(ctx, "orderUC.UpdateConfirm", "custom")
	defer span.End()

	statusID := 3

	return u.orderRepo.UpdateConfirm(ctx, id, statusID)
}
