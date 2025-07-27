package main

import (
	"log"
	"strings"
	"tctssf/config"
	"tctssf/controllers"
	//"tctssf/middleware"
	"tctssf/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	// Initialize database
	config.InitDB()
	defer config.CloseDB()
	
	app := fiber.New(fiber.Config{
		Prefork: false, // Set to true in production for better performance
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return ctx.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})
	
	// Middleware
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization, User-ID, User-Role",
		ExposeHeaders: "Authorization",
	}))
	
	// Add cache control for static files
	app.Use(func(c *fiber.Ctx) error {
		// For HTML files, disable caching during development
		if strings.HasSuffix(c.Path(), ".html") {
			c.Set("Cache-Control", "no-cache, no-store, must-revalidate")
			c.Set("Pragma", "no-cache")
			c.Set("Expires", "0")
		}
		return c.Next()
	})
	
	// Initialize controllers
	authController := controllers.NewAuthController()
	userController := controllers.NewUserController()
	adminController := controllers.NewAdminController()
	loanController := controllers.NewLoanController()
	treasurerController := controllers.NewTreasurerController()
	
	// Setup routes
	routes.SetupRoutes(app, authController, userController, adminController, loanController, treasurerController)
	
	log.Println("Server starting on port 3000...")
	log.Println("Frontend will be served from ./frontend directory")
	log.Println("Make sure your frontend files are placed in the ./frontend directory")
	log.Printf("Default Superadmin: superadmin@tctssf.rw / admin123")
	log.Printf("Default Admin: admin@tctssf.rw / admin123")
	log.Printf("Default Treasurer: treasurer@tctssf.rw / treasurer123")
	log.Fatal(app.Listen(":3000"))
}