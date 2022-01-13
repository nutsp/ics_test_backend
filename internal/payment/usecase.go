package payment

import (
	"app/internal/models"
	"context"
	"mime/multipart"
)

type UseCase interface {
	Create(ctx context.Context, payment *models.Payment) error
	UploadReceipt(ctx context.Context, mpf *multipart.Form) (string, error)
}
