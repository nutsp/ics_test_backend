package repository

import (
	"app/internal/models"
	"app/internal/order"
	"app/pkg/utils"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"go.elastic.co/apm"
)

type orderRepo struct {
	db *sqlx.DB
}

func NewOrderRepository(db *sqlx.DB) order.Repository {
	return &orderRepo{db: db}
}

func (r *orderRepo) checkProductInStock(ctx context.Context, productID, amount int64) (bool, error) {
	span, ctx := apm.StartSpan(ctx, "orderRepo.checkProductInStock", "custom")
	defer span.End()

	queryCount := `SELECT COUNT(*) FROM product 
		WHERE id = ? AND total_amount >= ?`
	var totalCount int = 0
	if err := r.db.GetContext(ctx, &totalCount, queryCount, productID, amount); err != nil {
		return false, err
	}

	if totalCount == 0 {
		return false, nil
	}

	return true, nil
}

func (r *orderRepo) Create(ctx context.Context, order *models.Order) error {
	span, ctx := apm.StartSpan(ctx, "orderRepo.Create", "custom")
	defer span.End()

	tx, err := r.db.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}

	createOrderStr := `INSERT INTO orders (user_id, total_price, address, status_id, create_at, update_at) 
		VALUES (?,?,?,?,?,?)`
	result, err := tx.ExecContext(ctx, createOrderStr, order.UserID, order.TotalPrice, order.Address, order.StatusID,
		order.CreateAt, order.UpdateAt)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	orderID, _ := result.LastInsertId()

	createOrderList := `INSERT INTO order_list (order_id, product_id, amount, price_per_unit)
		VALUES (?,?,?,?)`
	for _, od := range order.OrderList {
		inStock, err := r.checkProductInStock(ctx, od.ProductID, od.Amount)
		if err != nil {
			_ = tx.Rollback()
			return err
		}

		if inStock {
			_, err := tx.ExecContext(ctx, createOrderList, orderID, od.ProductID, od.Amount, od.PricePerUnit)
			if err != nil {
				_ = tx.Rollback()
				return err
			}

			updateProductStr := `UPDATE product SET total_amount = (total_amount - ?), update_at = ?
				WHERE id = ?`

			_, err = tx.ExecContext(ctx, updateProductStr, od.Amount, order.CreateAt, od.ProductID)
			if err != nil {
				_ = tx.Rollback()
				return err
			}
		} else {
			_ = tx.Rollback()
			return errors.New("not enough product")
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

	countOrderStr := `SELECT COUNT(*) FROM orders o
		LEFT JOIN order_status os ON o.status_id = os.id`
	orderListStr := `SELECT o.id, o.user_id, o.total_price, o.address, o.status_id, os.status, o.payment_id, o.create_at, o.update_at
		FROM orders o
		LEFT JOIN order_status os ON o.status_id = os.id`

	if filter.BeginDate != "" {
		sliceCon = append(sliceCon, ` From_unixtime(o.create_at, '%d-%m-%Y') >= '`+filter.BeginDate+`' `)
	}

	if filter.EndDate != "" {
		sliceCon = append(sliceCon, ` From_unixtime(o.create_at, '%d-%m-%Y') <= '`+filter.EndDate+`' `)
	}

	if filter.Status != "" {
		sliceCon = append(sliceCon, ` os.status = '`+filter.Status+`' `)
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

	countOrderStr := `SELECT COUNT(*) FROM orders o
		LEFT JOIN order_status os ON o.status_id = os.id`
	orderListStr := `SELECT o.id, o.user_id, o.total_price, o.address, o.status_id, os.status, o.payment_id, o.create_at, o.update_at
		FROM orders o
		LEFT JOIN order_status os ON o.status_id = os.id`

	sliceCon = append(sliceCon, ` o.user_id = '`+filter.UserID+`' `)
	if filter.BeginDate != "" {
		sliceCon = append(sliceCon, ` From_unixtime(o.create_at, '%d-%m-%Y') >= '`+filter.BeginDate+`' `)
	}

	if filter.EndDate != "" {
		sliceCon = append(sliceCon, ` From_unixtime(o.create_at, '%d-%m-%Y') <= '`+filter.EndDate+`' `)
	}

	if filter.Status != "" {
		sliceCon = append(sliceCon, ` os.status  = '`+filter.Status+`' `)
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

	var order = new(models.OrderBase)
	getOrderByID := `SELECT o.id, o.user_id, o.total_price, o.address, o.status_id, os.status, o.payment_id, o.create_at, o.update_at
		FROM orders o
		LEFT JOIN order_status os ON o.status_id = os.id
		WHERE o.id = ?`
	if err := r.db.QueryRowxContext(ctx, getOrderByID, id).StructScan(order); err != nil {
		return nil, err
	}

	countOrderList := `SELECT COUNT(*) FROM order_list o WHERE o.order_id = ?`
	var totalCount int
	if err := r.db.GetContext(ctx, &totalCount, countOrderList, id); err != nil {
		return nil, err
	}

	if totalCount == 0 {
		return order, nil
	}

	getOrderList := `SELECT p.id, p.name, p.gender, p.category_id, c.category, p.size_id, s.size,
		p.price_per_unit, o.amount
		FROM order_list o
		LEFT JOIN product p ON o.product_id = p.id
		LEFT JOIN category c ON c.id = p.category_id
		LEFT JOIN size s ON s.id = p.size_id
		WHERE o.order_id = ?`
	rows, err := r.db.QueryxContext(ctx, getOrderList, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orderList = make([]*models.OrderDetailBase, 0, totalCount)
	for rows.Next() {
		var orderDetail models.OrderDetailBase

		if err := rows.StructScan(&orderDetail); err != nil {
			return nil, err
		}
		orderList = append(orderList, &orderDetail)
	}

	order.OrderList = orderList
	return order, nil
}

func (r *orderRepo) UpdateCancel(ctx context.Context, id int, statusID int) error {
	span, ctx := apm.StartSpan(ctx, "orderRepo.UpdateCancel", "custom")
	defer span.End()

	updateOrderStr := `UPDATE orders SET status_id = ?, update_at = ? WHERE id = ?`
	_, err := r.db.ExecContext(ctx, updateOrderStr, statusID, time.Now().Unix(), id)
	if err != nil {
		return err
	}

	return nil
}

func (r *orderRepo) UpdateConfirm(ctx context.Context, id int, statusID int) error {
	span, ctx := apm.StartSpan(ctx, "orderRepo.UpdateConfirm", "custom")
	defer span.End()

	updateOrderStr := `UPDATE orders SET status_id = ?, update_at = ? WHERE id = ?`
	_, err := r.db.ExecContext(ctx, updateOrderStr, statusID, time.Now().Unix(), id)
	if err != nil {
		return err
	}

	return nil
}
