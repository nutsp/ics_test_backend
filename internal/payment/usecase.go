package payment

import (
	"app/internal/models"
	"context"
)

type UseCase interface {
	Create(ctx context.Context, payment *models.Payment) error
}
