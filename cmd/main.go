package main

import (
	"github.com/gin-gonic/gin"
	"github.com/kin-ark/GroAcademy/internal/database"
)

func main() {
  router := gin.Default()

  database.ConnectDB()

  router.Run()
}