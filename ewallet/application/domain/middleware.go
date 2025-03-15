package domain

import "github.com/gin-gonic/gin"

type Middleware interface {
	GetSessionCustomer() gin.HandlerFunc
}
