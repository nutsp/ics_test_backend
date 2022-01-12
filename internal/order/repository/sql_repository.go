package repository

import (
	"app/internal/models"
	"app/internal/order"
	"app/pkg/utils"
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"go.elastic.co/apm"
)

type orderRepo struct {
	db *sqlx.DB
}

func NewOrderRepository(db *sqlx.DB) order.Repository {
	return &orderRepo{db: db}
}

func (r *orderRepo) Create(ctx context.Context, order *models.Order) error {
	span, ctx := apm.StartSpan(ctx, "orderRepo.Create", "custom")
	defer span.End()

	tx, err := r.db.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}

	createOrderStr := `INSERT INTO orders (total_price, address, status, create_at, update_at) 
		VALUES (?,?,?,?,?)`
	result, err := tx.ExecContext(ctx, createOrderStr, order.TotalPrice, order.Address, order.Status,
		order.CreateAt, order.UpdateAt)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	orderID, _ := result.LastInsertId()

	createOrderList := `INSERT INTO order_list (order_id, product_id, amount, price_per_unit)
		VALUES (?,?,?,?)`
	for _, od := range order.OrderList {
		_, err := tx.ExecContext(ctx, createOrderList, orderID, od.ProductID, od.Amount, od.PricePerUnit)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (r *orderRepo) GetAll(ctx context.Context, filter *models.OrderFilter, query *utils.PaginationQuery) (*models.OrderList, error) {
	span, ctx := apm.StartSpan(ctx, "orderRepo.GetAll", "custom")
	defer span.End()

	return &models.OrderList{}, nil
}
