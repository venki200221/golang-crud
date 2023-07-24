package routes

import (
	"go_basics/controllers"
	"github.com/gofiber/fiber/v2"
)

func UserRoute(app *fiber.App){
app.Post("/user",controllers.CreateUser)
app.Get("/user/:userid",controllers.GetAUser)
app.Put("/user/:userid",controllers.EditUser)
app.Delete("/user/:userid",controllers.DeleteUser)
app.Get("/users",controllers.GetAllUsers)
}