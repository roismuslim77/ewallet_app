package payment

type MidtransPaymentBankTransferRequest struct {
	AmountTopUp   float64 `json:"amount_topup"`
	AmountService int     `json:"amount_service"`
	PaymentType   string  `json:"payment_type"`
	BankName      string  `json:"bank_name"`
	BankCode      string  `json:"bank_code"`
}

type MidtransWebhookRequest struct {
	VaNumbers []struct {
		VaNumber string `json:"va_number"`
		Bank     string `json:"bank"`
	} `json:"va_numbers"`
	TransactionTime   string `json:"transaction_time"`
	TransactionStatus string `json:"transaction_status"`
	TransactionId     string `json:"transaction_id"`
	StatusMessage     string `json:"status_message"`
	StatusCode        string `json:"status_code"`
	SignatureKey      string `json:"signature_key"`
	SettlementTime    string `json:"settlement_time"`
	PaymentType       string `json:"payment_type"`
	OrderId           string `json:"order_id"`
	GrossAmount       string `json:"gross_amount"`
	FraudStatus       string `json:"fraud_status"`
	ExpiryTime        string `json:"expiry_time"`
	Currency          string `json:"currency"`
}

type TransferOtherCustomer struct {
	Amount              float64 `json:"amount"`
	UserDestinationCode string  `json:"user_destination_code"`
}
