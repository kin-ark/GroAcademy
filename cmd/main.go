package main

import (
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

	router.Run()
}
