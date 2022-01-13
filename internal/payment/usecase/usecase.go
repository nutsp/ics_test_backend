package usecase

import (
	"app/internal/models"
	"app/internal/payment"
	"app/pkg/utils"
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.elastic.co/apm"
)

var allowedImagesContentType = map[string]string{
	"image/png":  "png",
	"image/jpg":  "jpg",
	"image/jpeg": "jpeg",
}

const (
	maxSize = 4194304 // 4194304 = 4MB
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

	statusID := 2 // Delivery
	now := time.Now().Unix()
	payment.CreateAt = now

	return u.paymentRepo.Create(ctx, payment, statusID)
}

func (u *paymentUC) UploadReceipt(ctx context.Context, mpf *multipart.Form) (string, error) {
	span, ctx := apm.StartSpan(ctx, "paymentUC.UploadReceipt", "custom")
	defer span.End()

	var c *fiber.Ctx

	now := time.Now()
	files := mpf.File["receipt"]
	if len(files) < 0 && len(files) > 1 {
		return "", errors.New("must have files min 1 or max 4")
	}

	sizeFiles := 0
	for i := 0; i < len(files); i++ {
		contentTypes := files[i].Header.Get("Content-Type")
		if len(contentTypes) < 1 {
			return "", errors.New("not allowed header content-type")
		}

		_, allowed := allowedImagesContentType[contentTypes]
		if !allowed {
			return "", errors.New("not allowed header content-type")
		}
		sizeFiles = sizeFiles + int(files[i].Size)
	}

	if sizeFiles > maxSize {
		return "", errors.New("file size not over 4MB")
	}

	dir := fmt.Sprintf("storage/receipt/%s/", now.Format("2006/01/02"))
	if err := utils.CheckDir("./" + dir); err != nil {
		return "", errors.New("create dir err: " + err.Error())
	}

	contentTypes := files[0].Header.Get("Content-Type")
	fileType := allowedImagesContentType[contentTypes]
	newDir := fmt.Sprintf("/%s%d.%s", dir, now.Unix(), fileType)
	c.SaveFile(files[0], fmt.Sprintf(".%s", newDir))

	return newDir, nil
}
