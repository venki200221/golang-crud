package main

import (
	"go_basics/configs"
	"go_basics/routes"
	"github.com/gofiber/fiber/v2"
	
)

func main(){
	app:=fiber.New()
	configs.ConnectDB()
	routes.UserRoute(app)
	app.Get("/",func(c *fiber.Ctx)error {
		return c.JSON(&fiber.Map{"data":"hello world"})
	})
	app.Listen(":3000")
}