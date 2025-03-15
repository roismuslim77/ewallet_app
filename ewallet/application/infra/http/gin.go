package infrahttp

import (
	"fmt"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	midtrans "simple-go/application/adapter/midtrans"
	"simple-go/application/domain/auth"
	"simple-go/application/domain/healthcheck"
	"simple-go/application/domain/payment"
)

type Router struct {
	router     *gin.Engine
	port       string
	db         *gorm.DB
	middleware *Middleware
	adapter    Adapter
}

func handleCors() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Writer.Header().Set("Access-Control-Allow-Origin", ctx.Request.Header.Get("Origin"))
		ctx.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		ctx.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-Client-ID")
		ctx.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if ctx.Request.Method == "OPTIONS" {
			ctx.AbortWithStatus(204)
			return
		}

		ctx.Next()
	}
}

func NewRouter(port string, pg *gorm.DB) *Router {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.Use(gzip.Gzip(gzip.DefaultCompression))
	router.Use(handleCors())

	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "ok",
		})
	})

	return &Router{
		router: router,
		port:   port,
		db:     pg,
	}
}

func (r *Router) Run() {
	fmt.Println("server running at port", r.port)

	baseHealthCheck := r.router.Group("/health-check")
	baseAuth := r.router.Group("/auth")
	baseWallet := r.router.Group("wallet")

	r.BuildHealthCheck(baseHealthCheck)
	r.BuildAuth(baseAuth)
	r.BuildPayment(baseWallet)

	r.router.Run(fmt.Sprintf(":%s", r.port))
}

func (r *Router) BuildHealthCheck(router *gin.RouterGroup) {
	hc := healthcheck.NewRouterHttp(router, r.db)
	hc.RegisterRoute()
}
func (r *Router) BuildAuth(router *gin.RouterGroup) {
	auth := auth.NewRouterHttp(router, r.db, r.middleware)
	auth.RegisterRoute()
}
func (r *Router) BuildPayment(router *gin.RouterGroup) {
	payment := payment.NewRouterHttp(router, r.db, r.adapter, r.middleware)
	payment.RegisterRoute()
}

func (r *Router) SetAdapter() *Router {
	midtrans := midtrans.NewMidtransAdapter()
	buildAdapter := Adapter{
		midtrans: midtrans,
	}

	adapter := NewAdapter().Build(buildAdapter)

	r.adapter = adapter

	return r
}

func (r *Router) SetMiddleware(db *gorm.DB) *Router {
	mid := NewBuilderMiddleware()
	r.middleware = mid

	return r
}
