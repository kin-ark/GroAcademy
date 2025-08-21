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

	router.LoadHTMLGlob("internal/templates/*")

	database.ConnectDB()

	routes.RegisterFEoutes(router)
	routes.RegisterRoutes(router)

	router.Run()
}
