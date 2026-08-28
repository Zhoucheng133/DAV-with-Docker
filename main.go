package main

import (
	"log"

	"dav_docker/routes"
	"dav_docker/utils"

	"github.com/gofiber/fiber/v2"
)

func main() {
	if err := utils.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	if err := utils.InitEnv(); err != nil {
		log.Fatalf("Failed to initialize environment secrets: %v", err)
	}

	utils.InitRunningWebDAVServers()

	app := fiber.New()

	app.Get("/api/init", routes.HandleInit)
	app.Post("/api/register", routes.HandleRegister)
	app.Post("/api/login", routes.HandleLogin)
	app.Get("/api/refresh", routes.HandleRefresh)

	app.Get("/api/auth", routes.AuthMiddleware, routes.HandleAuth)
	app.Post("/api/logout", routes.AuthMiddleware, routes.HandleLogout)

	app.Get("/api/config/list", routes.AuthMiddleware, routes.HandleConfigList)
	app.Post("/api/config/add", routes.AuthMiddleware, routes.HandleConfigAdd)
	app.Post("/api/config/edit/:id", routes.AuthMiddleware, routes.HandleConfigEdit)
	app.Delete("/api/config/del/:id", routes.AuthMiddleware, routes.HandleConfigDel)
	app.Post("/api/config/run/:id", routes.AuthMiddleware, routes.HandleConfigRun)
	app.Post("/api/config/stop/:id", routes.AuthMiddleware, routes.HandleConfigStop)

	log.Println("Server starting on :3000...")
	log.Fatal(app.Listen(":3000"))
}
