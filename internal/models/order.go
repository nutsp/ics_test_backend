package models

type Order struct {
	UserID     string         `json:"-" form:"-"`
	TotalPrice float64        `json:"total_price" form:"-"`
	Address    string         `json:"address" form:"address"`
	Status     int            `json:"status" form:"-"`
	OrderList  []*OrderDetail `json:"order_list" form:"order_list"`
	CreateAt   int64          `json:"create_at" form:"-"`
	UpdateAt   int64          `json:"update_at" form:"-"`
}

type OrderDetail struct {
	ProductID    int64   `json:"product_id" form:"product_id"`
	Amount       int64   `json:"amount" form:"amount"`
	PricePerUnit float64 `json:"price_per_unit" form:"price_per_unit"`
}

type OrderBase struct {
	ID         int            `json:"id" db:"id"`
	UserID     string         `json:"-" db:"user_id"`
	TotalPrice float64        `json:"total_price" db:"total_price"`
	Address    string         `json:"address" db:"address"`
	Status     int            `json:"status" db:"status"`
	PaymentID  int            `json:"payment_id" db:"payment_id"`
	OrderList  []*OrderDetail `json:"order_list" db:"order_list"`
	CreateAt   int64          `json:"create_at" db:"create_at"`
	UpdateAt   int64          `json:"update_at" db:"update_at"`
}

type OrderFilter struct {
	UserID    string `json:"user_id" query:"user_id"`
	BeginDate string `json:"begin_date" query:"begin_date"`
	EndDate   string `json:"end_date" query:"end_date"`
	Status    string `json:"status" query:"status"`
}

type OrderList struct {
	TotalCount int          `json:"total_count"`
	TotalPages int          `json:"total_pages"`
	Page       int          `json:"page"`
	Size       int          `json:"limit"`
	HasMore    bool         `json:"has_more"`
	Orders     []*OrderBase `json:"orders"`
}
