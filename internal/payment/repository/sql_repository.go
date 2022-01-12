package repository

import (
	"app/internal/models"
	"app/internal/payment"
	"context"

	"github.com/jmoiron/sqlx"
	"go.elastic.co/apm"
)

type paymentRepo struct {
	db sqlx.DB
}

func NewPaymentRepository(db sqlx.DB) payment.Repository {
	return &paymentRepo{db: db}
}

func (r *paymentRepo) Create(ctx context.Context, payment *models.Payment) error {
	span, ctx := apm.StartSpan(ctx, "paymentRepo.Create", "custom")
	defer span.End()

	return nil
}
