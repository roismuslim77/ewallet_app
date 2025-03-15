package auth

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

func NewRouterHttp(router *gin.RouterGroup, db *gorm.DB, middle domain.Middleware) domain.HttpHandler {
	repository := NewRepository(db)
	service := NewService(repository)

	handler := NewHandler(&service)

	return &RouterHttp{
		router:     router,
		handler:    handler,
		middleware: middle,
	}
}

func (r RouterHttp) RegisterRoute() {
	r.router.POST("/register", r.handler.RegisterCustomer)
	r.router.POST("/login", r.handler.LoginCustomer)

}
