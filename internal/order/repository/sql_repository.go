package repository

import (
	"app/internal/models"
	"app/internal/order"
	"app/pkg/utils"
	"context"
	"database/sql"
	"fmt"

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

	createOrderStr := `INSERT INTO orders (user_id, total_price, address, status, create_at, update_at) 
		VALUES (?,?,?,?,?,?)`
	result, err := tx.ExecContext(ctx, createOrderStr, order.UserID, order.TotalPrice, order.Address, order.Status,
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

	var sliceCon []string
	var queryCount, queryAll string

	countOrderStr := `SELECT COUNT(*) FROM orders o`
	orderListStr := `SELECT o.id, o.user_id, o.total_price, o.address, o.status, o.payment_id, o.create_at, o.update_at
		FROM orders o`

	if filter.BeginDate != "" {
		sliceCon = append(sliceCon, ` From_unixtime(o.create_at, '%d-%m-%Y') >= '`+filter.BeginDate+`' `)
	}

	if filter.EndDate != "" {
		sliceCon = append(sliceCon, ` From_unixtime(o.create_at, '%d-%m-%Y') <= '`+filter.EndDate+`' `)
	}

	if filter.Status != "" {
		sliceCon = append(sliceCon, ` o.status  = '`+filter.Status+`' `)
	}

	condition := utils.ConcatCondition(sliceCon)
	if condition != "" {
		queryCount = fmt.Sprintf("%s WHERE %s", countOrderStr, condition)
		queryAll = fmt.Sprintf("%s WHERE %s", orderListStr, condition)
	} else {
		queryCount = fmt.Sprintf("%s", countOrderStr)
		queryAll = fmt.Sprintf("%s", orderListStr)
	}

	var totalCount int = 0
	if err := r.db.GetContext(ctx, &totalCount, queryCount); err != nil {
		return nil, err
	}

	if totalCount == 0 {
		return &models.OrderList{
			TotalCount: totalCount,
			TotalPages: utils.GetTotalPages(totalCount, query.GetSize()),
			Page:       query.GetPage(),
			Size:       query.GetSize(),
			HasMore:    utils.GetHasMore(query.GetPage(), totalCount, query.GetSize()),
			Orders:     make([]*models.OrderBase, 0),
		}, nil
	}

	queryAll += ` LIMIT ? OFFSET ?`
	rows, err := r.db.QueryxContext(ctx, queryAll, query.GetLimit(), query.GetOffset())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orderList = make([]*models.OrderBase, 0, query.GetSize())
	for rows.Next() {
		var order models.OrderBase
		if err := rows.StructScan(&order); err != nil {
			return nil, err
		}
		orderList = append(orderList, &order)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &models.OrderList{
		TotalCount: totalCount,
		TotalPages: utils.GetTotalPages(totalCount, query.GetSize()),
		Page:       query.GetPage(),
		Size:       query.GetSize(),
		HasMore:    utils.GetHasMore(query.GetPage(), totalCount, query.GetSize()),
		Orders:     orderList,
	}, nil
}

func (r *orderRepo) GetByUserID(ctx context.Context, filter *models.OrderFilter, query *utils.PaginationQuery) (*models.OrderList, error) {
	span, ctx := apm.StartSpan(ctx, "orderRepo.GetByUserID", "custom")
	defer span.End()

	var sliceCon []string
	var queryCount, queryAll string

	countOrderStr := `SELECT COUNT(*) FROM orders o`
	orderListStr := `SELECT o.id, o.user_id, o.total_price, o.address, o.status, o.payment_id, o.create_at, o.update_at
		FROM orders o`

	sliceCon = append(sliceCon, ` o.user_id = '`+filter.UserID+`' `)
	if filter.BeginDate != "" {
		sliceCon = append(sliceCon, ` From_unixtime(o.create_at, '%d-%m-%Y') >= '`+filter.BeginDate+`' `)
	}

	if filter.EndDate != "" {
		sliceCon = append(sliceCon, ` From_unixtime(o.create_at, '%d-%m-%Y') <= '`+filter.EndDate+`' `)
	}

	if filter.Status != "" {
		sliceCon = append(sliceCon, ` o.status  = '`+filter.Status+`' `)
	}

	condition := utils.ConcatCondition(sliceCon)
	if condition != "" {
		queryCount = fmt.Sprintf("%s WHERE %s", countOrderStr, condition)
		queryAll = fmt.Sprintf("%s WHERE %s", orderListStr, condition)
	} else {
		queryCount = fmt.Sprintf("%s", countOrderStr)
		queryAll = fmt.Sprintf("%s", orderListStr)
	}

	var totalCount int = 0
	if err := r.db.GetContext(ctx, &totalCount, queryCount); err != nil {
		return nil, err
	}

	if totalCount == 0 {
		return &models.OrderList{
			TotalCount: totalCount,
			TotalPages: utils.GetTotalPages(totalCount, query.GetSize()),
			Page:       query.GetPage(),
			Size:       query.GetSize(),
			HasMore:    utils.GetHasMore(query.GetPage(), totalCount, query.GetSize()),
			Orders:     make([]*models.OrderBase, 0),
		}, nil
	}

	queryAll += ` LIMIT ? OFFSET ?`
	rows, err := r.db.QueryxContext(ctx, queryAll, query.GetLimit(), query.GetOffset())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orderList = make([]*models.OrderBase, 0, query.GetSize())
	for rows.Next() {
		var order models.OrderBase
		if err := rows.StructScan(&order); err != nil {
			return nil, err
		}
		orderList = append(orderList, &order)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &models.OrderList{
		TotalCount: totalCount,
		TotalPages: utils.GetTotalPages(totalCount, query.GetSize()),
		Page:       query.GetPage(),
		Size:       query.GetSize(),
		HasMore:    utils.GetHasMore(query.GetPage(), totalCount, query.GetSize()),
		Orders:     orderList,
	}, nil
}

func (r *orderRepo) GetByID(ctx context.Context, id int) (*models.OrderBase, error) {
	span, ctx := apm.StartSpan(ctx, "orderRepo.GetByID", "custom")
	defer span.End()

	return nil, nil
}
