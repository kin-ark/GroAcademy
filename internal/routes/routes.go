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
	userService := services.NewUserService(userRepo)
	userController := controllers.NewUserController(userService)

	courseRepo := repositories.NewCourseRepository()
	courseService := services.NewCourseService(courseRepo)
	courseController := controllers.NewCourseController(courseService)

	moduleRepo := repositories.NewModuleRepository()
	moduleService := services.NewModuleService(moduleRepo, courseRepo)
	moduleController := controllers.NewModuleController(moduleService)

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
		courses := api.Group("/courses")
		courses.Use(middlewares.RequireAuth)
		{
			courses.POST("", middlewares.RequireAdmin, courseController.PostCourse)
			courses.GET("", courseController.GetAllCourses)
			courses.GET("/:id", courseController.GetCourseByID)
			courses.PUT("/:id", middlewares.RequireAdmin, courseController.PutCourse)
			courses.DELETE("/:id", middlewares.RequireAdmin, courseController.DeleteCourseByID)
			courses.POST("/:id/buy", courseController.BuyCourse)
			courses.GET("/my-courses", courseController.GetMyCourses)

			// CRUD Module
			courses.POST("/:id/modules", middlewares.RequireAdmin, moduleController.PostModule)
			courses.GET("/:id/modules", moduleController.GetModules)
			courses.PATCH("/:id/modules/reorder", middlewares.RequireAdmin, moduleController.ReorderModules)
			courses.PATCH("/:id/modules/complete", moduleController.MarkModuleAsComplete)
		}

		// CRUD Module Route
		modules := api.Group("/modules")
		modules.Use(middlewares.RequireAuth)
		{
			modules.GET("/:id", moduleController.GetModuleById)
			modules.PUT("/:id", middlewares.RequireAdmin, moduleController.PutModule)
			modules.DELETE("/:id", middlewares.RequireAdmin, moduleController.DeleteModuleByID)
		}

		// CRUD User Route
		api.GET("/users", userController.GetUsers)
		api.GET("/users/:id", userController.GetUserById)
		// api.POST("/users/:id/balance", controllers.PostUserBalance)
		// api.PUT("/users/:id", controllers.PutUserById)
		// api.DELETE("/users/:id", controllers.DeleteUserById)
	}
}
