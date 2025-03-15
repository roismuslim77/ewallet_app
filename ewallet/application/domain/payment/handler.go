package payment

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"log"
	"net/http"
	"simple-go/pkg/response"
)

type Service interface {
	CustomerPaymentTopUp(ctx context.Context, customerId int, req MidtransPaymentBankTransferRequest) (string, response.ErrorResponse)
	WebhookPaymentTopUp(ctx context.Context, req MidtransWebhookRequest) response.ErrorResponse
	CustomerWalletTrx(ctx context.Context, customerId int) ([]CustomerWalletHistory, response.ErrorResponse)
	CustomerTransferOther(ctx context.Context, customerId int, req TransferOtherCustomer) response.ErrorResponse
}

type handler struct {
	service Service
}

func NewHandler(svc Service) handler {
	return handler{
		service: svc,
	}
}

func (h handler) CustomerPaymentTopUp(ctx *gin.Context) {
	customerId := ctx.GetInt("customerId")
	var req MidtransPaymentBankTransferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Println(err.Error())
		resp := response.Error("22102").WithStatusCode(http.StatusBadRequest)

		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			resp.WithArgsMessage(ve[0].Field(), ve[0].Tag())
		}

		ctx.AbortWithStatusJSON(resp.StatusCode, resp)
		return
	}

	vaNumber, err := h.service.CustomerPaymentTopUp(ctx, customerId, req)
	if !err.IsNoError {
		resp := response.Error(err.Code).WithError(err.Message).WithStatusCode(err.StatusCode)
		ctx.AbortWithStatusJSON(resp.StatusCode, resp)
		return
	}

	resp := response.Success("22151").WithData(gin.H{"va_number": vaNumber})
	ctx.JSON(resp.StatusCode, resp)
}

func (h handler) WebhookPaymentTopUp(ctx *gin.Context) {
	var req MidtransWebhookRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Println(err.Error())
		resp := response.Error("22102").WithStatusCode(http.StatusBadRequest)

		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			resp.WithArgsMessage(ve[0].Field(), ve[0].Tag())
		}

		ctx.AbortWithStatusJSON(resp.StatusCode, resp)
		return
	}

	err := h.service.WebhookPaymentTopUp(ctx, req)
	if !err.IsNoError {
		resp := response.Error(err.Code).WithError(err.Message).WithStatusCode(err.StatusCode)
		ctx.AbortWithStatusJSON(resp.StatusCode, resp)
		return
	}

	resp := response.Success("22154")
	ctx.JSON(resp.StatusCode, resp)
}

func (h handler) CustomerWalletTrx(ctx *gin.Context) {
	customerId := ctx.GetInt("customerId")
	trx, err := h.service.CustomerWalletTrx(ctx, customerId)
	if !err.IsNoError {
		resp := response.Error(err.Code).WithError(err.Message).WithStatusCode(err.StatusCode)
		ctx.AbortWithStatusJSON(resp.StatusCode, resp)
		return
	}

	resp := response.Success("22154").WithData(trx)
	ctx.JSON(resp.StatusCode, resp)
}

func (h handler) CustomerTransferOther(ctx *gin.Context) {
	customerId := ctx.GetInt("customerId")
	var req TransferOtherCustomer
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Println(err.Error())
		resp := response.Error("22102").WithStatusCode(http.StatusBadRequest)

		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			resp.WithArgsMessage(ve[0].Field(), ve[0].Tag())
		}

		ctx.AbortWithStatusJSON(resp.StatusCode, resp)
		return
	}

	err := h.service.CustomerTransferOther(ctx, customerId, req)
	if !err.IsNoError {
		resp := response.Error(err.Code).WithError(err.Message).WithStatusCode(err.StatusCode)
		ctx.AbortWithStatusJSON(resp.StatusCode, resp)
		return
	}

	resp := response.Success("22157")
	ctx.JSON(resp.StatusCode, resp)
}
