package models_backup

type Outlets struct {
	ID         uint64  `json:"id" gorm:"column:id;primary_key;autoIncrement"`
	MerchantId uint64  `json:"merchant_id" gorm:"column:merchant_id"`
	OutletName string  `json:"outlet_name" gorm:"column:outlet_name"`
	CreatedAt  string  `json:"created_at" gorm:"column:created_at"`
	CreatedBy  *uint64 `json:"created_by" gorm:"column:created_by"`
	UpdatedAt  string  `json:"updated_at" gorm:"column:updated_at"`
	UpdatedBy  *uint64 `json:"updated_by" gorm:"column:updated_by"`
}
