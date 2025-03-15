package payment

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"simple-go/application/domain"
)

type RouterHttp struct {
	router     *gin.RouterGroup
	handler    handler
	middleware domain.Middleware
}

func NewRouterHttp(router *gin.RouterGroup, db *gorm.DB, adapter domain.Adapter, middle domain.Middleware) domain.HttpHandler {
	adapterData := domain.NewAdapter().Build(adapter)

	repository := NewRepository(db)
	service := NewService(repository, adapterData)

	handler := NewHandler(&service)

	return &RouterHttp{
		router:     router,
		handler:    handler,
		middleware: middle,
	}
}

func (r RouterHttp) RegisterRoute() {
	r.router.POST("/topup", r.middleware.GetSessionCustomer(), r.handler.CustomerPaymentTopUp)
	r.router.GET("/history", r.middleware.GetSessionCustomer(), r.handler.CustomerWalletTrx)
	r.router.POST("/pay", r.middleware.GetSessionCustomer(), r.handler.CustomerTransferOther)

	r.router.POST("/webhook", r.handler.WebhookPaymentTopUp)
}
