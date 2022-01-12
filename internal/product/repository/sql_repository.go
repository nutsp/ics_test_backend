package repository

import (
	"app/internal/models"
	"app/internal/product"
	"app/pkg/utils"
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"go.elastic.co/apm"
)

type productRepo struct {
	db *sqlx.DB
}

func NewProductRepository(db *sqlx.DB) product.Repository {
	return &productRepo{db: db}
}

func (r *productRepo) Create(ctx context.Context, product *models.Product) error {
	span, ctx := apm.StartSpan(ctx, "productRepo.Create", "custom")
	defer span.End()

	tx, err := r.db.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}

	createProductStr := `INSERT INTO product (name, gender, category_id, size_id,
		total_amount, price_per_unit, create_at, update_at)
		VALUES (?,?,?,?,?,?,?,?)`
	result, err := tx.ExecContext(ctx, createProductStr, product.Name, product.Gender, product.CategoryID,
		product.SizeID, product.Amount, product.PricePerUnit, product.CreateAt, product.UpdateAt)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	product.ID, _ = result.LastInsertId()

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (r *productRepo) GetAll(ctx context.Context, filter *models.ProductFilter, query *utils.PaginationQuery) (*models.ProductList, error) {
	span, ctx := apm.StartSpan(ctx, "productRepo.GetAll", "custom")
	defer span.End()

	var sliceCon []string
	var queryCount, queryAll string

	countProductStr := `SELECT COUNT(*) FROM product p
		LEFT JOIN category c ON c.id = p.category_id
		LEFT JOIN size s ON s.id = p.size_id`
	productListStr := `SELECT p.id, p.name, p.gender, p.category_id, c.category, p.size_id, s.size,
		p.total_amount, p.price_per_unit, p.create_at, p.update_at 
		FROM product p
		LEFT JOIN category c ON c.id = p.category_id
		LEFT JOIN size s ON s.id = p.size_id`

	sliceCon = append(sliceCon, ` p.total_amount > 0 `)

	if filter.Gender != "" {
		sliceCon = append(sliceCon, ` UPPER(p.gender) = '`+strings.ToUpper(filter.Gender)+`' `)
	}

	if filter.Category != "" {
		sliceCon = append(sliceCon, ` UPPER(c.category) = '`+strings.ToUpper(filter.Category)+`' `)
	}

	if filter.Size != "" {
		sliceCon = append(sliceCon, ` UPPER(s.size) = '`+strings.ToUpper(filter.Size)+`' `)
	}

	condition := utils.ConcatCondition(sliceCon)
	if condition != "" {
		queryCount = fmt.Sprintf("%s WHERE %s", countProductStr, condition)
		queryAll = fmt.Sprintf("%s WHERE %s", productListStr, condition)
	} else {
		queryCount = fmt.Sprintf("%s", countProductStr)
		queryAll = fmt.Sprintf("%s", productListStr)
	}

	var totalCount int = 0
	if err := r.db.GetContext(ctx, &totalCount, queryCount); err != nil {
		return nil, err
	}

	if totalCount == 0 {
		return &models.ProductList{
			TotalCount: totalCount,
			TotalPages: utils.GetTotalPages(totalCount, query.GetSize()),
			Page:       query.GetPage(),
			Size:       query.GetSize(),
			HasMore:    utils.GetHasMore(query.GetPage(), totalCount, query.GetSize()),
			Products:   make([]*models.ProductBase, 0),
		}, nil
	}

	queryAll += ` LIMIT ? OFFSET ?`
	rows, err := r.db.QueryxContext(ctx, queryAll, query.GetLimit(), query.GetOffset())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var producutList = make([]*models.ProductBase, 0, query.GetSize())
	for rows.Next() {
		var product models.ProductBase
		if err := rows.StructScan(&product); err != nil {
			return nil, err
		}
		producutList = append(producutList, &product)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &models.ProductList{
		TotalCount: totalCount,
		TotalPages: utils.GetTotalPages(totalCount, query.GetSize()),
		Page:       query.GetPage(),
		Size:       query.GetSize(),
		HasMore:    utils.GetHasMore(query.GetPage(), totalCount, query.GetSize()),
		Products:   producutList,
	}, nil
}
