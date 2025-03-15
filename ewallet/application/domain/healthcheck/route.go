package healthcheck

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"simple-go/application/domain"
)

type RouterHttp struct {
	router  *gin.RouterGroup
	handler handler
}

func NewRouterHttp(router *gin.RouterGroup, db *gorm.DB) domain.HttpHandler {
	handler := NewHandler(db)

	return &RouterHttp{
		router:  router,
		handler: handler,
	}
}

func (r RouterHttp) RegisterRoute() {
	r.router.GET("/", r.handler.Healthcheck)
}
