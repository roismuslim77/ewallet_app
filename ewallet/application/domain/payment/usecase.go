package payment

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"simple-go/application/config"
	"simple-go/application/domain"
	"simple-go/application/entity"
	"simple-go/helper"
	"simple-go/pkg/response"
	"time"
)

type Repository interface {
	GetCustomerById(ctx context.Context, id int) (entity.Customer, error)
	GetCustomerByCode(ctx context.Context, code string) (entity.Customer, error)

	GetWalletTrx(ctx context.Context, customerId int) ([]entity.WalletTransaction, error)
	GetWalletTrxByReferenceId(ctx context.Context, referenceId string) (entity.WalletTransaction, error)
	CreateWalletTransaction(ctx context.Context, req entity.WalletTransaction) (entity.WalletTransaction, error)
	UpdateWalletTransaction(ctx context.Context, req entity.WalletTransaction, id int) (entity.WalletTransaction, error)

	GetWalletCustomer(ctx context.Context, customerId int) (entity.Wallet, error)
	CreateWalletCustomer(ctx context.Context, req entity.Wallet) (entity.Wallet, error)
	UpdateWalletCustomer(ctx context.Context, req entity.Wallet, id int) (entity.Wallet, error)
}

type service struct {
	repository Repository
	adapter    domain.AdapterData
}

func NewService(repo Repository, adapter domain.AdapterData) service {
	return service{
		repository: repo,
		adapter:    adapter,
	}
}

func (s service) WebhookPaymentTopUp(ctx context.Context, req MidtransWebhookRequest) response.ErrorResponse {
	trx, err := s.repository.GetWalletTrxByReferenceId(ctx, req.OrderId)
	if err != nil {
		return *response.Error("22101").WithError(err.Error()).WithStatusCode(http.StatusBadRequest)
	}
	if trx.Status {
		return *response.Error("22101").WithError("status has been paid").WithStatusCode(http.StatusBadRequest)
	}

	hash := helper.SHA512Hash(req.OrderId + req.StatusCode + req.GrossAmount + config.GetString(config.CFG_MIDTRANS_SERVER_KEY, ""))
	if hash != req.SignatureKey {
		return *response.Error("22101").WithError("signature key not match").WithStatusCode(http.StatusBadRequest)
	}

	updateWalletTrx := entity.WalletTransaction{
		Status:      true,
		Description: fmt.Sprintf(`Topup Saldo Wallet sebesar %v berhasil`, trx.Amount),
		UpdatedAt:   time.Now(),
	}
	_, err = s.repository.UpdateWalletTransaction(ctx, updateWalletTrx, trx.ID)
	if err != nil {
		return *response.Error("22101").WithError(err.Error()).WithStatusCode(http.StatusInternalServerError)
	}

	customerWallet, err := s.repository.GetWalletCustomer(ctx, trx.CustomerId)
	if err != nil {
		return *response.Error("22101").WithError(err.Error()).WithStatusCode(http.StatusInternalServerError)
	}

	if customerWallet.IsEmpty {
		//create new wallet
		createWallet := entity.Wallet{
			CustomerId: trx.CustomerId,
			Status:     true,
			Balance:    trx.Amount,
			CreatedAt:  time.Now(),
		}
		_, err := s.repository.CreateWalletCustomer(ctx, createWallet)
		if err != nil {
			return *response.Error("22101").WithError(err.Error()).WithStatusCode(http.StatusInternalServerError)
		}
	} else {
		//update wallet
		updateWallet := entity.Wallet{
			Balance: customerWallet.Balance + trx.Amount,
		}
		_, err := s.repository.UpdateWalletCustomer(ctx, updateWallet, customerWallet.ID)
		if err != nil {
			return *response.Error("22101").WithError(err.Error()).WithStatusCode(http.StatusInternalServerError)
		}
	}

	return *response.NotError()
}

func (s service) CustomerPaymentTopUp(ctx context.Context, customerId int, req MidtransPaymentBankTransferRequest) (string, response.ErrorResponse) {
	orderId := fmt.Sprintf(`TPP-0%v`, time.Now().UnixNano())

	//get customer data
	customer, err := s.repository.GetCustomerById(ctx, customerId)
	if err != nil {
		return "", *response.Error("22101").WithError(err.Error()).WithStatusCode(http.StatusInternalServerError)
	}
	//create data wallet transaction
	dataWalletTrx := entity.WalletTransaction{
		CustomerId:   customerId,
		ReferenceId:  orderId,
		Amount:       req.AmountTopUp,
		MutationType: "topup",
		Category:     "transaction",
		BankName:     req.BankName,
		Description:  fmt.Sprintf(`Topup Saldo Wallet sebesar %v sedang di proses.`, req.AmountTopUp),
		Status:       false,
		CreatedAt:    time.Now(),
	}
	respDataWalletTrx, err := s.repository.CreateWalletTransaction(ctx, dataWalletTrx)
	if err != nil {
		return "", *response.Error("22101").WithError(err.Error()).WithStatusCode(http.StatusInternalServerError)
	}

	//generate va number from midtrans
	encryptByte, err := json.Marshal(gin.H{
		"payment_type": req.PaymentType,
		"transaction_details": gin.H{
			"order_id":     orderId,
			"gross_amount": req.AmountTopUp,
		},
		"bank_transfer": gin.H{
			"bank": req.BankCode,
		},
		"customer_details": gin.H{
			"last_name": customer.FullName,
		},
	})
	if err != nil {
		return "", *response.Error("22101").WithError(err.Error()).WithStatusCode(http.StatusInternalServerError)
	}

	requestVAData, err := s.adapter.Midtrans.RequestVirtualAccount(encryptByte)
	if err != nil {
		return "", *response.Error("22101").WithError(err.Error()).WithStatusCode(http.StatusInternalServerError)
	}

	//update wallet trx
	updateWalletTrx := entity.WalletTransaction{NoAcc: requestVAData.VaNumbers[0].VaNumber}
	_, err = s.repository.UpdateWalletTransaction(ctx, updateWalletTrx, respDataWalletTrx.ID)
	if err != nil {
		return "", *response.Error("22101").WithError(err.Error()).WithStatusCode(http.StatusInternalServerError)
	}

	return requestVAData.VaNumbers[0].VaNumber, *response.NotError()
}

func (s service) CustomerWalletTrx(ctx context.Context, customerId int) ([]CustomerWalletHistory, response.ErrorResponse) {
	var data []CustomerWalletHistory
	trx, err := s.repository.GetWalletTrx(ctx, customerId)
	if err != nil {
		return data, *response.Error("22101").WithError(err.Error()).WithStatusCode(http.StatusInternalServerError)
	}

	for _, transaction := range trx {
		temp := CustomerWalletHistory{
			Id:           transaction.ID,
			ReferenceId:  transaction.ReferenceId,
			Amount:       transaction.Amount,
			MutationType: transaction.MutationType,
			Category:     transaction.Category,
			BankName:     transaction.BankName,
			NoAcc:        transaction.NoAcc,
			Description:  transaction.Description,
			Status:       transaction.Status,
			CreatedAt:    transaction.CreatedAt,
			UpdatedAt:    transaction.UpdatedAt,
		}

		data = append(data, temp)
	}

	return data, *response.NotError()
}

func (s service) CustomerTransferOther(ctx context.Context, customerId int, req TransferOtherCustomer) response.ErrorResponse {
	//check customer id with destination code
	customer, err := s.repository.GetCustomerById(ctx, customerId)
	if err != nil {
		return *response.Error("22101").WithError(err.Error()).WithStatusCode(http.StatusInternalServerError)
	}
	if customer.IsEmpty {
		return *response.Error("22101").WithError("not found").WithStatusCode(http.StatusBadRequest)
	}
	if customer.Code == req.UserDestinationCode {
		return *response.Error("22101").WithError("can't transfer to own account").WithStatusCode(http.StatusBadRequest)
	}

	customerDest, err := s.repository.GetCustomerByCode(ctx, req.UserDestinationCode)
	if err != nil {
		return *response.Error("22101").WithError(err.Error()).WithStatusCode(http.StatusInternalServerError)
	}
	if customerDest.IsEmpty {
		return *response.Error("22101").WithError("customer destination not found").WithStatusCode(http.StatusBadRequest)
	}

	//check current balance
	balanceFrom, err := s.repository.GetWalletCustomer(ctx, customer.ID)
	if err != nil {
		return *response.Error("22101").WithError(err.Error()).WithStatusCode(http.StatusInternalServerError)
	}
	if balanceFrom.IsEmpty {
		return *response.Error("22101").WithError("balance not found").WithStatusCode(http.StatusBadRequest)
	}
	if balanceFrom.Balance < req.Amount {
		return *response.Error("22101").WithError("balance wallet not enough").WithStatusCode(http.StatusBadRequest)
	}

	balanceDest, err := s.repository.GetWalletCustomer(ctx, customerDest.ID)
	if err != nil {
		return *response.Error("22101").WithError(err.Error()).WithStatusCode(http.StatusInternalServerError)
	}
	if balanceDest.IsEmpty {
		create := entity.Wallet{
			CustomerId: customerDest.ID,
			Status:     true,
			Balance:    req.Amount,
			CreatedAt:  time.Now(),
		}
		_, err := s.repository.CreateWalletCustomer(ctx, create)
		if err != nil {
			return *response.Error("22101").WithError(err.Error()).WithStatusCode(http.StatusInternalServerError)
		}
	} else {
		//update customer balance destination
		updateWallet := entity.Wallet{Balance: balanceDest.Balance + req.Amount, UpdatedAt: time.Now()}
		_, err = s.repository.UpdateWalletCustomer(ctx, updateWallet, balanceDest.ID)
		if err != nil {
			return *response.Error("22101").WithError(err.Error()).WithStatusCode(http.StatusInternalServerError)
		}
	}
	//update customer balance source
	updateWallet := entity.Wallet{Balance: balanceFrom.Balance - req.Amount, UpdatedAt: time.Now()}
	_, err = s.repository.UpdateWalletCustomer(ctx, updateWallet, balanceFrom.ID)
	if err != nil {
		return *response.Error("22101").WithError(err.Error()).WithStatusCode(http.StatusInternalServerError)
	}

	//create data wallet transaction destination
	orderId := fmt.Sprintf(`%v%v-%v`, customer.Code, customerDest.Code, time.Now().UnixNano())
	dataWalletTrx := entity.WalletTransaction{
		CustomerId:   customerDest.ID,
		ReferenceId:  orderId,
		Amount:       req.Amount,
		MutationType: "topup",
		Category:     "transfer",
		Description:  fmt.Sprintf(`Transfer Saldo Wallet sebesar %v dari %v berhasil diterima.`, req.Amount, customer.FullName),
		Status:       true,
		CreatedAt:    time.Now(),
	}
	_, err = s.repository.CreateWalletTransaction(ctx, dataWalletTrx)
	if err != nil {
		return *response.Error("22101").WithError(err.Error()).WithStatusCode(http.StatusInternalServerError)
	}

	//create data wallet transaction source
	dataWalletTrx = entity.WalletTransaction{
		CustomerId:   customer.ID,
		ReferenceId:  orderId,
		Amount:       req.Amount,
		MutationType: "withdraw",
		Category:     "transfer",
		Description:  fmt.Sprintf(`Transfer Saldo Wallet sebesar %v ke %v berhasil.`, req.Amount, customerDest.FullName),
		Status:       true,
		CreatedAt:    time.Now(),
	}
	_, err = s.repository.CreateWalletTransaction(ctx, dataWalletTrx)
	if err != nil {
		return *response.Error("22101").WithError(err.Error()).WithStatusCode(http.StatusInternalServerError)
	}

	return *response.NotError()
}
