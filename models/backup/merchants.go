package models_backup

type Merchants struct {
	ID           uint64  `json:"id" gorm:"column:id;primary_key;autoIncrement"`
	UserId       uint64  `json:"user_id" gorm:"column:user_id"`
	MerchantName string  `json:"merchant_name" gorm:"column:merchant_name"`
	CreatedAt    string  `json:"created_at" gorm:"column:created_at"`
	CreatedBy    *uint64 `json:"created_by" gorm:"column:created_by"`
	UpdatedAt    string  `json:"updated_at" gorm:"column:updated_at"`
	UpdatedBy    *uint64 `json:"updated_by" gorm:"column:updated_by"`
}
