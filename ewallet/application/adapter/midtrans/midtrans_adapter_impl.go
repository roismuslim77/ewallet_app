package adapter

import (
	"encoding/base64"
	"encoding/json"
	"simple-go/application/config"
	"simple-go/pkg"
)

type midtransAdapter struct {
	client      pkg.HttpClient
	defaultPath string
	apiKey      string
}

func (m midtransAdapter) RequestVirtualAccount(req []byte) (MidtransPaymentBankTransferResponse, error) {
	var respData MidtransPaymentBankTransferResponse
	encoded := base64.StdEncoding.EncodeToString([]byte(m.apiKey))

	path := "/charge"
	headers := pkg.RequestHeader{
		ContentType:   "application/json",
		Authorization: "Basic " + encoded,
	}

	resp, err := m.client.PostThirdParty(headers, m.defaultPath+path, req)
	if err != nil {
		return respData, err
	}

	err = json.Unmarshal(resp, &respData)
	if err != nil {
		return respData, err
	}

	return respData, nil
}

func NewMidtransAdapter() MidtransAdapter {
	rpiMidtrans := config.GetString(config.CFG_MIDTRANS_URL, "")
	rpiMidtransKey := config.GetString(config.CFG_MIDTRANS_SERVER_KEY, "")
	client := pkg.NewHttpClient()

	return midtransAdapter{
		client:      client,
		defaultPath: rpiMidtrans,
		apiKey:      rpiMidtransKey,
	}
}
