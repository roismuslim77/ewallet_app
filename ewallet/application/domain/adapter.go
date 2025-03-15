package domain

import midtrans "simple-go/application/adapter/midtrans"

type AdapterData struct {
	Midtrans midtrans.MidtransAdapter
}

func NewAdapter() AdapterData {
	return AdapterData{}
}

func (a AdapterData) Build(adapter Adapter) AdapterData {
	return AdapterData{
		Midtrans: adapter.GetMidtransAdapter(),
	}
}

type Adapter interface {
	GetMidtransAdapter() midtrans.MidtransAdapter
}
