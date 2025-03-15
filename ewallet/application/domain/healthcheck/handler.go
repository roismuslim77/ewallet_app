package healthcheck

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"simple-go/pkg/response"
)

type handler struct {
	db *gorm.DB
}

func NewHandler(db *gorm.DB) handler {
	return handler{
		db: db,
	}
}

func (h handler) Healthcheck(ctx *gin.Context) {
	var postgres string

	postgresDb, err := h.db.DB()
	if err != nil {
		fmt.Println(err)
		postgres = "down"
	}

	postgresErr := postgresDb.Ping()
	if postgresErr != nil {
		fmt.Println(postgresErr)
		postgres = "down"
	} else {
		postgres = "up"
	}

	resp := response.Success("22152")
	ctx.JSON(resp.StatusCode, gin.H{
		"response": "up",
		"databases": gin.H{
			"postgres": postgres,
		},
	})
}
