package payment

import "time"

type CustomerWalletHistory struct {
	Id            int       `json:"id"`
	CustomerId    int       `json:"customer_id"`
	ReferenceId   string    `json:"reference_id"`
	InvoiceNumber string    `json:"invoice_number"`
	Amount        float64   `json:"amount"`
	MutationType  string    `json:"mutation_type"`
	Category      string    `json:"category"`
	BankName      string    `json:"bank_name"`
	NoAcc         string    `json:"no_acc"`
	Description   string    `json:"description"`
	Note          string    `json:"note"`
	Status        bool      `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
