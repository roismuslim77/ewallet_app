package infra

import (
	"gorm.io/gorm"
	infrahttp "simple-go/application/infra/http"
)

type Infra interface {
	Run()
}

type InfraBuilder struct {
}

func NewInfraFactory() *InfraBuilder {
	return &InfraBuilder{}
}

func (i *InfraBuilder) CreateInfraHttp(port string, pg *gorm.DB) (Infra, error) {
	return infrahttp.NewRouter(port, pg).SetAdapter().SetMiddleware(pg), nil
}
