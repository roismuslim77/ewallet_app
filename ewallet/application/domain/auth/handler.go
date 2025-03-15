package auth

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
	RegisterCustomer(ctx context.Context, req RegisterUserRequest) response.ErrorResponse
	LoginCustomer(ctx context.Context, req LoginCustomerReq) (string, response.ErrorResponse)
}

type handler struct {
	service Service
}

func NewHandler(svc Service) handler {
	return handler{
		service: svc,
	}
}

func (h handler) RegisterCustomer(ctx *gin.Context) {
	var req RegisterUserRequest
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

	err := h.service.RegisterCustomer(ctx, req)
	if !err.IsNoError {
		resp := response.Error(err.Code).WithError(err.Message).WithStatusCode(err.StatusCode)
		ctx.AbortWithStatusJSON(resp.StatusCode, resp)
		return
	}

	resp := response.Success("22155")
	ctx.JSON(resp.StatusCode, resp)
}

func (h handler) LoginCustomer(ctx *gin.Context) {
	var req LoginCustomerReq
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

	token, err := h.service.LoginCustomer(ctx, req)
	if !err.IsNoError {
		resp := response.Error(err.Code).WithError(err.Message).WithStatusCode(err.StatusCode)
		ctx.AbortWithStatusJSON(resp.StatusCode, resp)
		return
	}

	resp := response.Success("22156").WithData(gin.H{"token": token})
	ctx.JSON(resp.StatusCode, resp)
}
