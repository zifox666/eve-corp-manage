package main

import (
	"eve-corp-manager/config"
	"eve-corp-manager/initialize"
	"eve-corp-manager/router"
)

func main() {
	config.InitConfig()

	initialize.StartUp()

	port := config.AppConfig.App.Port

	if port == "" {
		port = "8005"
	}

	if err := router.InitRouters(":" + port); err != nil {
		panic(err)
	}

}
