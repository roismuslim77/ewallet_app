package infrahttp

import midtrans "simple-go/application/adapter/midtrans"

type Adapter struct {
	midtrans midtrans.MidtransAdapter
}

func NewAdapter() Adapter {
	return Adapter{}
}

func (a Adapter) GetMidtransAdapter() midtrans.MidtransAdapter {
	return a.midtrans
}

func (a Adapter) Build(adapter Adapter) Adapter {
	return adapter
}
