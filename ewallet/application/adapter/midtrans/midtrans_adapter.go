package adapter

type MidtransAdapter interface {
	RequestVirtualAccount(req []byte) (MidtransPaymentBankTransferResponse, error)
}
