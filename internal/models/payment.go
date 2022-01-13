package models

type Payment struct {
	OrderID     int64   `json:"order_id" form:"order_id"`
	BankID      int64   `json:"bank_id" form:"bank_id"`
	Price       float64 `json:"price" form:"price"`
	ReceiptPath string  `json:"receipt_path" form:"receipt_path"`
	CreateAt    int64   `json:"create_at" form:"-"`
}

type PaymentBase struct {
}
