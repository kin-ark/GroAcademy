package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/kin-ark/GroAcademy/internal/controllers"
	"github.com/kin-ark/GroAcademy/internal/middlewares"
	"github.com/kin-ark/GroAcademy/internal/repositories"
	"github.com/kin-ark/GroAcademy/internal/services"
)

func RegisterRoutes(r *gin.Engine) {
	// FOR TESTING PURPOSE ONLY, WILL BE DELETED LATER
	userRepo := repositories.NewUserRepository()
	authService := services.NewAuthService(userRepo)
	authController := controllers.NewAuthController(authService)

	courseRepo := repositories.NewCourseRepository()
	courseService := services.NewCourseService(courseRepo)
	courseController := controllers.NewCourseController(courseService)

	api := r.Group("/api")
	{
		// Auth Route
		auth := api.Group("/auth")
		{
			auth.POST("/login", authController.Login)
			auth.POST("/register", authController.Register)
			auth.GET("/self", middlewares.RequireAuth, authController.GetSelf)
		}

		// CRUD Course Route
		api.POST("/courses", courseController.PostCourse)
		api.GET("/courses", courseController.GetAllCourses)
		api.GET("/courses/:id", courseController.GetCourseByID)
		api.PUT("/courses/:id", courseController.PutCourse)
		// api.DELETE("/courses/:id", controllers.DeleteCourseByID)
		// api.POST("/courses/:id/buy", controllers.BuyCourse)
		// api.GET("/courses/my-courses", controllers.GetUserCourse)

		// CRUD Module Route
		// api.POST("/courses/:courseId/modules", controllers.PostModules)
		// api.GET("/courses/:courseId/modules", controllers.GetModulesByCourse)
		// api.GET("/modules/:id", controllers.GetModuleById)
		// api.PUT("/modules/:id", controllers.PutModuleById)
		// api.DELETE("/modules/:id", controllers.DeleteModuleById)
		// api.PATCH("/courses/:courseId/modules/reorder", controllers.ReorderCourseModules)
		// api.PATCH("/courses/:courseId/modules/complete", controllers.MarkModuleAsComplete)

		// CRUD User Route
		// api.GET("/users", controllers.GetUsers)
		// api.GET("/users/:id", controllers.GetUserById)
		// api.POST("/users/:id/balance", controllers.PostUserBalance)
		// api.PUT("/users/:id", controllers.PutUserById)
		// api.DELETE("/users/:id", controllers.DeleteUserById)
	}
}
