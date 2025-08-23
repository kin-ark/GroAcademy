package routes

import (
	"log"
	"html/template"

	"github.com/gin-gonic/gin"
	"github.com/kin-ark/GroAcademy/internal/controllers"
	"github.com/kin-ark/GroAcademy/internal/middlewares"
	"github.com/kin-ark/GroAcademy/internal/repositories"
	"github.com/kin-ark/GroAcademy/internal/services"
)

func SetupHTMLRenderer(router *gin.Engine) {
	funcMap := template.FuncMap{
		"add": func(a, b int) int { return a + b },
		"sub": func(a, b int) int { return a - b },
		"mul": func(a, b int) int { return a * b },
	}

	tmpl := template.New("").Funcs(funcMap)

	tmpl = template.Must(tmpl.ParseGlob("internal/templates/components/*.html"))
	tmpl = template.Must(tmpl.ParseGlob("internal/templates/*.html"))

	router.SetHTMLTemplate(tmpl)

	log.Println("Loaded HTML Templates:")
	for _, t := range tmpl.Templates() {
		log.Println(t.Name())
	}
}

func RegisterFEoutes(r *gin.Engine) {
	userRepo := repositories.NewUserRepository()
	courseRepo := repositories.NewCourseRepository()
	moduleRepo := repositories.NewModuleRepository()

	userService := services.NewUserService(userRepo)
	courseService := services.NewCourseService(courseRepo)
	moduleService := services.NewModuleService(moduleRepo, courseRepo)

	fc := controllers.NewFEController(userService, courseService, moduleService)

	r.GET("/login", fc.ShowLoginPage)

	r.GET("/register", fc.ShowRegisterPage)

	r.GET("/courses", middlewares.FERequireAuth, fc.GetCoursesPage)

	r.GET("/course/:id", middlewares.FERequireAuth, fc.GetCourseDetailPage)
}
