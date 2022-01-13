package models

type Product struct {
	ID           int64   `json:"id" form:"-"`
	Name         string  `json:"name" form:"name"`
	Gender       string  `json:"gender" form:"gender"`
	CategoryID   int64   `json:"category_id" form:"category_id"`
	SizeID       int64   `json:"size_id" form:"size_id"`
	Amount       int64   `json:"amount" form:"amount"`
	PricePerUnit float64 `json:"price_per_unit" form:"price_per_unit"`
	CreateAt     int64   `json:"create_at" form:"-"`
	UpdateAt     int64   `json:"update_at" form:"-"`
}

type ProductBase struct {
	ID           int64   `json:"id" db:"id"`
	Name         string  `json:"name" db:"name"`
	Gender       string  `json:"gender" db:"gender"`
	CategoryID   int64   `json:"-" db:"category_id"`
	Category     string  `json:"category" db:"category"`
	SizeID       int64   `json:"-" db:"size_id"`
	Size         string  `json:"size" db:"size"`
	TotalAmount  int64   `json:"total_amount" db:"total_amount"`
	PricePerUnit float64 `json:"price_per_unit" db:"price_per_unit"`
	CreateAt     int64   `json:"create_at" db:"create_at"`
	UpdateAt     int64   `json:"update_at" db:"update_at"`
}

type ProductFilter struct {
	Gender   string `json:"gender" query:"gender"`
	Category string `json:"category_id" query:"category"`
	Size     string `json:"size_id" query:"size"`
}

type ProductList struct {
	TotalCount int            `json:"total_count"`
	TotalPages int            `json:"total_pages"`
	Page       int            `json:"page"`
	Size       int            `json:"limit"`
	HasMore    bool           `json:"has_more"`
	Products   []*ProductBase `json:"products"`
}
