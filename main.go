package main

import (
	"demariot-backend/controllers"
	"demariot-backend/database"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	database.InitDB()

	app := fiber.New()

	app.Use(logger.New(logger.Config{
		Format: "${time} ${ip} ${method} ${path} ${status}\n",
	}))

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000",
		AllowMethods:     "GET,POST, PUT, DELETE",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
	}))

	app.Static("/uploads", "./uploads")

	app.Post("/auth/login", controllers.Login)
	app.Post("/api/register", controllers.Register)
	app.Get("/api/users", controllers.GetUsers)
	app.Get("/api/users/:id", controllers.GetUserById)
	app.Put("/api/users/:id/role", controllers.UpdateUserRole)
	app.Put("/users/:id", controllers.UpdateUser)
	app.Put("/api/users/:id/change_password", controllers.ChangePassword)
	app.Post("/user/:id/profile_picture", controllers.UploadProfilePicture)
	app.Delete("/api/users/:id", controllers.DeleteUser)

	fmt.Println("Server running on port 3001")
	app.Listen(":3001")
}
