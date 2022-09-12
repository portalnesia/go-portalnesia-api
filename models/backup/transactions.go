package models_backup

type Transactions struct {
	ID         uint64  `json:"id" gorm:"column:id;primary_key;autoIncrement"`
	MerchantId uint64  `json:"merchant_id" gorm:"column:merchant_id"`
	OutletId   uint64  `json:"outlet_id" gorm:"column:outlet_id"`
	BillTotal  float64 `json:"bill_total" gorm:"column:bill_total"`
	CreatedAt  string  `json:"created_at" gorm:"column:created_at"`
	CreatedBy  *uint64 `json:"created_by" gorm:"column:created_by"`
	UpdatedAt  string  `json:"updated_at" gorm:"column:updated_at"`
	UpdatedBy  *uint64 `json:"updated_by" gorm:"column:updated_by"`
}
