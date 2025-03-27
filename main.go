package main

import (
	"services-auth/config"
	"services-auth/controller"
	"services-auth/middlewares"

	"github.com/labstack/echo/v4"
)

func main() {

	e := echo.New()
	config.DatabaseInit()

	auth := e.Group("api/v1/auth")

	auth.POST("/register", controller.Register)
	auth.POST("/login", controller.Login)
	auth.GET("/user", middlewares.Auth(controller.Home))

	e.Logger.Fatal(e.Start(":6060"))
}
