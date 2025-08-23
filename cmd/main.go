package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/kin-ark/GroAcademy/internal/database"
	"github.com/kin-ark/GroAcademy/internal/routes"
)

func main() {
	router := gin.Default()

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

	router.Run(":" + port)
}
