package entity

import "time"

type CustomerAuth struct {
	ID           int    `gorm:"column:id;type:int;primaryKey;autoIncrement:true;unique" json:"id"`
	CustomerId   int    `gorm:"column:customer_id;type:int;unique" json:"customer_id"`
	Password     string `gorm:"column:password;type:string;size:255" json:"password"`
	Token        string `gorm:"column:token;type:string;size:255" json:"token"`
	RefreshToken string `gorm:"column:refresh_token;type:string;size:255" json:"refresh_token"`

	CreatedAt time.Time `gorm:"column:created_at;" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;" json:"updated_at"`

	Customer Customer `json:"customer" gorm:"foreignKey:CustomerId;references:ID"`

	// custom
	IsEmpty bool `gorm:"-" json:"-"`
}

func (t CustomerAuth) TableName() string {
	return "customer_auth"
}
