package routes

import (
	"todo/controllers"
	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App){
	app.Get("/home",controllers.Home)
	app.Get("/reg",controllers.RegisterTMPL)
	app.Get("/log",controllers.LoginTMPL)
	app.Get("/",controllers.Index)
	app.Post("/api/addtask", controllers.AddTask)
	app.Post("/api/register", controllers.Register)
	app.Post("/api/login", controllers.Login)
	app.Get("/api/user", controllers.User)
	app.Get("/api/gettask/:id", controllers.GetTasks)
	app.Delete("/api/deletetask/:taskid", controllers.DeleteTask)
	app.Post("/api/logout", controllers.Logout)
	
}