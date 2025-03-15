package entity

import "time"

type Wallet struct {
	ID         int     `gorm:"column:id;type:int;primaryKey;autoIncrement:true;unique" json:"id"`
	CustomerId int     `gorm:"column:customer_id;type:int;unique" json:"customer_id"`
	Status     bool    `gorm:"column:status;type:boolean;" json:"status"`
	Balance    float64 `gorm:"column:balance;type:numeric;" json:"balance"`

	CreatedAt time.Time `gorm:"column:created_at;" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;" json:"updated_at"`

	Customer Customer `json:"customer" gorm:"foreignKey:CustomerId;references:ID"`

	// custom
	IsEmpty bool `gorm:"-" json:"-"`
}

func (t Wallet) TableName() string {
	return "wallets"
}
