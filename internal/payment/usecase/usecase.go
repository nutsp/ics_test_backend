package usecase

import (
	"app/internal/models"
	"app/internal/payment"
	"context"

	"go.elastic.co/apm"
)

type paymentUC struct {
	paymentRepo payment.Repository
}

func NewPaymentUseCase(paymentRepo payment.Repository) payment.UseCase {
	return &paymentUC{paymentRepo: paymentRepo}
}

func (u *paymentUC) Create(ctx context.Context, payment *models.Payment) error {
	span, ctx := apm.StartSpan(ctx, "paymentUC.Create", "custom")
	defer span.End()

	return u.paymentRepo.Create(ctx, payment)
}
