package models

type Order struct {
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
}

type OrderFilter struct {
}

type OrderList struct {
	TotalCount int          `json:"total_count"`
	TotalPages int          `json:"total_pages"`
	Page       int          `json:"page"`
	Size       int          `json:"limit"`
	HasMore    bool         `json:"has_more"`
	Orders     []*OrderBase `json:"orders"`
}
