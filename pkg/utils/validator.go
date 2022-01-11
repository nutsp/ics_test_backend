package utils

import (
	"github.com/go-playground/validator/v10"
)

// Use a single instance of Validate, it caches struct info
var validate *validator.Validate

func init() {
	validate = validator.New()
}

// Validate struct fields
// func ValidateStruct(ctx context.Context, s interface{}) error {
// 	return validate.StructCtx(ctx, s)
// }

// ValidateStruct ...
func ValidateStruct(s interface{}) error {
	validate := validator.New()
	err := validate.Struct(s)
	if err != nil {
		return err
	}
	return nil
}
