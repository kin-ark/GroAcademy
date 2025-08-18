package main

import (
	"github.com/gin-gonic/gin"
	"github.com/kin-ark/GroAcademy/internal/database"
	"github.com/kin-ark/GroAcademy/internal/routes"
)

func main() {
	router := gin.Default()

	database.ConnectDB()

	routes.RegisterRoutes(router)

	router.Run()
}
