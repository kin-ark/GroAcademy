package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/kin-ark/GroAcademy/internal/controllers"
	"github.com/kin-ark/GroAcademy/internal/middlewares"
	"github.com/kin-ark/GroAcademy/internal/repositories"
	"github.com/kin-ark/GroAcademy/internal/services"
)

func RegisterRoutes(r *gin.Engine) {
	// Dependency Injection
	userRepo := repositories.NewUserRepository()
	courseRepo := repositories.NewCourseRepository()
	moduleRepo := repositories.NewModuleRepository()

	authService := services.NewAuthService(userRepo)
	userService := services.NewUserService(userRepo)
	courseService := services.NewCourseService(courseRepo)
	moduleService := services.NewModuleService(moduleRepo, courseRepo)

	authController := controllers.NewAuthController(authService)
	userController := controllers.NewUserController(userService)
	courseController := controllers.NewCourseController(courseService)
	moduleController := controllers.NewModuleController(moduleService)

	api := r.Group("/api")
	{
		registerAuthRoutes(api, &authController)
		registerCourseRoutes(api, &courseController, &moduleController)
		registerModuleRoutes(api, &moduleController)
		registerUserRoutes(api, &userController)
	}
}

func registerAuthRoutes(api *gin.RouterGroup, authController *controllers.AuthController) {
	auth := api.Group("/auth")
	{
		auth.POST("/login", authController.Login)
		auth.POST("/register", authController.Register)
		auth.GET("/self", middlewares.RequireAuth, authController.GetSelf)
	}
}

func registerCourseRoutes(api *gin.RouterGroup, courseController *controllers.CourseController, moduleController *controllers.ModuleController) {
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

		courses.POST("/:id/modules", middlewares.RequireAdmin, moduleController.PostModule)
		courses.GET("/:id/modules", moduleController.GetModules)
		courses.PATCH("/:id/modules/reorder", middlewares.RequireAdmin, moduleController.ReorderModules)
	}
}

func registerModuleRoutes(api *gin.RouterGroup, moduleController *controllers.ModuleController) {
	modules := api.Group("/modules")
	modules.Use(middlewares.RequireAuth)
	{
		modules.GET("/:id", moduleController.GetModuleById)
		modules.PUT("/:id", middlewares.RequireAdmin, moduleController.PutModule)
		modules.DELETE("/:id", middlewares.RequireAdmin, moduleController.DeleteModuleByID)
		modules.PATCH("/:id/complete", moduleController.MarkModuleAsComplete)
	}
}

func registerUserRoutes(api *gin.RouterGroup, userController *controllers.UserController) {
	users := api.Group("/users")
	users.Use(middlewares.RequireAuth, middlewares.RequireAdmin)
	{
		users.GET("", userController.GetUsers)
		users.GET("/:id", userController.GetUserById)
		users.POST("/:id/balance", userController.AddUserBalance)
		users.PUT("/:id", userController.PutUser)
		users.DELETE("/:id", userController.DeleteUserByID)
	}
}
