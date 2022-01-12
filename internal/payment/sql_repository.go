package payment

import (
	"app/internal/models"
	"context"
)

type Repository interface {
	Create(ctx context.Context, payment *models.Payment) error
}
