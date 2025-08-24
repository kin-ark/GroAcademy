package main

import (
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/kin-ark/GroAcademy/internal/database"
	"github.com/kin-ark/GroAcademy/internal/routes"
)

func main() {
	router := gin.Default()

	config := cors.DefaultConfig()

	allowedOrigins := os.Getenv("ALLOWED_ORIGINS")
	if allowedOrigins != "" {
		config.AllowOrigins = strings.Split(allowedOrigins, ",")
	} else {
		config.AllowOrigins = []string{
			"http://localhost:3000",
			"http://localhost:5173",
			"http://127.0.0.1:3000",
			"http://127.0.0.1:5173",
		}
	}

	config.AllowMethods = []string{
		"GET",
		"POST",
		"PUT",
		"DELETE",
		"OPTIONS",
		"PATCH",
	}

	config.AllowHeaders = []string{
		"Origin",
		"Content-Type",
		"Accept",
		"Authorization",
		"X-Requested-With",
		"X-CSRF-Token",
	}

	config.AllowCredentials = true

	config.MaxAge = 12 * 60 * 60

	router.Use(cors.New(config))

	router.Static("/static", "./static")
	router.Static("/uploads", "./uploads")

	database.ConnectDB()

	routes.SetupHTMLRenderer(router)
	routes.RegisterFEoutes(router)
	routes.RegisterRoutes(router)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// seeder := seeds.NewSeeder(database.DB)
	// err := seeder.SeedAll()
	// if err != nil {
	// 	log.Println(err.Error())
	// }

	router.Run(":" + port)
}
