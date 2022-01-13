package repository

import (
	"app/internal/models"
	"app/internal/payment"
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"go.elastic.co/apm"
)

type paymentRepo struct {
	db *sqlx.DB
}

func NewPaymentRepository(db *sqlx.DB) payment.Repository {
	return &paymentRepo{db: db}
}

func (r *paymentRepo) Create(ctx context.Context, payment *models.Payment, statusID int) error {
	span, ctx := apm.StartSpan(ctx, "paymentRepo.Create", "custom")
	defer span.End()

	tx, err := r.db.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}

	createPaymentStr := `INSERT INTO payment (bank_id, price, receipt_path, create_at)
		VALUES (?,?,?,?)`
	result, err := tx.ExecContext(ctx, createPaymentStr, payment.BankID, payment.Price,
		payment.ReceiptPath, payment.CreateAt)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	payID, _ := result.LastInsertId()
	updateOrderStr := `UPDATE orders SET status_id = ?, payment_id = ?, update_at = ? WHERE id = ?`
	_, err = tx.ExecContext(ctx, updateOrderStr, statusID, payID, payment.CreateAt, payment.OrderID)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}
