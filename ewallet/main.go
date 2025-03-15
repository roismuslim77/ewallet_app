package main

import (
	"simple-go/application/config"
	"simple-go/application/database"
	"simple-go/application/infra"
	"time"
)

func main() {
	jakartaTimeZone, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		panic(err)
	}

	time.Local = jakartaTimeZone

	err = config.LoadConfig(".env")
	if err != nil {
		panic(err)
	}

	pg, err := database.ConnectPostgres()
	if err != nil {
		panic(err)
	}

	infra := infra.NewInfraFactory()
	infraHttp, err := infra.CreateInfraHttp(config.GetString(config.CFG_APP_PORT, "3000"), pg)
	if err != nil {
		panic(err)
	}
	infraHttp.Run()
}
