package entity

import "time"

type WalletTransaction struct {
	ID            int     `gorm:"column:id;type:int;primaryKey;autoIncrement:true;unique" json:"id"`
	CustomerId    int     `gorm:"column:customer_id;type:int" json:"customer_id"`
	ReferenceId   string  `gorm:"column:reference_id;type:string;size:255" json:"reference_id"`
	InvoiceNumber string  `gorm:"column:invoice_number;type:string;size:255" json:"invoice_number"`
	Amount        float64 `gorm:"column:amount;type:numeric;" json:"amount"`
	MutationType  string  `gorm:"column:mutation_type;type:string;size:100" json:"mutation_type"`
	Category      string  `gorm:"column:category;type:string;size:100" json:"category"`
	BankName      string  `gorm:"column:bank_name;type:string;size:100" json:"bank_name"`
	NoAcc         string  `gorm:"column:no_acc;type:string;size:100" json:"no_acc"`
	Description   string  `gorm:"column:description;type:string;size:100" json:"description"`
	Notes         string  `gorm:"column:note;type:string;size:255" json:"note"`

	Status bool `gorm:"column:status;type:boolean;" json:"status"`

	CreatedAt time.Time `gorm:"column:created_at;" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;" json:"updated_at"`

	Customer Customer `json:"customer" gorm:"foreignKey:CustomerId;references:ID"`

	// custom
	IsEmpty bool `gorm:"-" json:"-"`
}

func (t Wallet) WalletTransaction() string {
	return "wallet_transactions"
}
