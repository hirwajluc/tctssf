package main

import (
	"log"
	"strings"
	"tctssf/config"
	"tctssf/controllers"
	_ "tctssf/docs" // Import generated docs
	"tctssf/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

// @title RDB SLS API
// @version 1.0
// @description RDB Staff Loan System - Savings and Loan Management System API
// @contact.name API Support
// @contact.email support@rdbsls.rw
// @license.name MIT
// @host localhost:3000
// @BasePath /api
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize structured logger
	config.InitLogger(cfg.Environment)
	defer config.CloseLogger()

	// Initialize database
	config.InitDB()
	defer config.CloseDB()

	// Initialize Redis (optional - will fall back to in-memory if unavailable)
	config.InitRedis(cfg.RedisURL)
	defer config.CloseRedis()

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
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${method} ${path} (${latency})\n",
	}))
	app.Use(recover.New())

	// CORS with configurable origins
	app.Use(cors.New(cors.Config{
		AllowOrigins:  cfg.AllowedOrigins,
		AllowMethods:  "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:  "Origin, Content-Type, Accept, Authorization, User-ID, User-Role",
		ExposeHeaders: "Authorization",
	}))

	// Rate limiting
	app.Use(limiter.New(limiter.Config{
		Max:        100, // 100 requests
		Expiration: 60,  // per minute
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(429).JSON(fiber.Map{
				"error": "Too many requests, please try again later",
			})
		},
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
	
	// Swagger documentation (MUST be before routes to avoid being caught by static handler)
	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	// Initialize controllers
	authController := controllers.NewAuthController()
	userController := controllers.NewUserController()
	adminController := controllers.NewAdminController()
	loanController := controllers.NewLoanController()
	treasurerController := controllers.NewTreasurerController()

	// Setup routes
	routes.SetupRoutes(app, authController, userController, adminController, loanController, treasurerController)

	// Server info
	serverAddr := ":" + cfg.ServerPort
	log.Printf("===========================================")
	log.Printf("RDB SLS Server Starting")
	log.Printf("===========================================")
	log.Printf("Environment: %s", cfg.Environment)
	log.Printf("Server: http://%s:%s", cfg.ServerHost, cfg.ServerPort)
	log.Printf("Database: %s@%s:%s/%s", cfg.DBUser, cfg.DBHost, cfg.DBPort, cfg.DBName)
	if config.RedisClient != nil {
		log.Printf("Redis: Connected ✓")
	} else {
		log.Printf("Redis: Not connected (using in-memory sessions)")
	}
	log.Printf("CORS Origins: %s", cfg.AllowedOrigins)
	log.Printf("Rate Limiting: 100 requests/minute")
	log.Printf("Frontend: ./frontend directory")
	log.Printf("-------------------------------------------")
	log.Printf("Default Credentials:")
	log.Printf("  Superadmin: superadmin@rdbsls.rw / admin123")
	log.Printf("  Admin: admin@rdbsls.rw / admin123")
	log.Printf("  Treasurer: treasurer@rdbsls.rw / treasurer123")
	log.Printf("===========================================")

	// Start server with optional TLS
	if cfg.EnableTLS {
		log.Printf("Starting HTTPS server on %s", serverAddr)
		log.Fatal(app.ListenTLS(serverAddr, cfg.TLSCertFile, cfg.TLSKeyFile))
	} else {
		log.Printf("Starting HTTP server on %s", serverAddr)
		log.Fatal(app.Listen(serverAddr))
	}
}