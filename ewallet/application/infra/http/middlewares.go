package infrahttp

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
	"log"
	"net/http"
	"simple-go/application/domain/auth"
	"simple-go/application/entity"
	"simple-go/helper"
	"simple-go/pkg/response"
)

type Middleware struct {
}

func NewBuilderMiddleware() *Middleware {
	return &Middleware{}
}

func (m *Middleware) AddHeader() gin.HandlerFunc {
	log.Println("header")
	return func(ctx *gin.Context) {
		ctx.Header("Content-Type", "application/json")
		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.Next()
	}
}

func (m Middleware) GetSessionCustomer() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var authHeader auth.CustomerAuthHeaderReq

		err := ctx.ShouldBindHeader(&authHeader)
		if err != nil {
			log.Println(err.Error(), " :err")
			resp := response.Error("22102")

			var ve validator.ValidationErrors
			if errors.As(err, &ve) {
				resp.WithArgsMessage(ve[0].Field(), ve[0].Tag())
			}

			ctx.AbortWithStatusJSON(resp.StatusCode, resp)
			return
		}

		tokenString := authHeader.Authorization
		log.Println(tokenString, "token")
		claims := &entity.Claims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(helper.GetJWTKey()), nil
		})
		if err != nil || !token.Valid {
			log.Println(err.Error())
			resp := response.Error("22150")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, resp)
			return
		}

		ctx.Set("customerId", claims.CustomerId)
		ctx.Next()
	}
}
