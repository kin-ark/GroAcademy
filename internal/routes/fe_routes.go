package routes

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kin-ark/GroAcademy/internal/controllers"
	"github.com/kin-ark/GroAcademy/internal/middlewares"
	"github.com/kin-ark/GroAcademy/internal/models"
	"github.com/kin-ark/GroAcademy/internal/repositories"
	"github.com/kin-ark/GroAcademy/internal/services"
)

func SetupHTMLRenderer(router *gin.Engine) {
	funcMap := template.FuncMap{
		"add":                 func(a, b int) int { return a + b },
		"sub":                 func(a, b int) int { return a - b },
		"mul":                 func(a, b int) int { return a * b },
		"moduleURL":           moduleURL,
		"moduleURLWithType":   moduleURLWithType,
		"moduleCompletionURL": moduleCompletionURL,
		"moduleTypeLabel":     moduleTypeLabel,
		"hasMultipleContent":  hasMultipleContent,
		"shouldShowPDF":       shouldShowPDF,
		"shouldShowVideo":     shouldShowVideo,
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

	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/login")
	})

	r.GET("/login", middlewares.RedirectIfAuthenticated, fc.ShowLoginPage)

	r.GET("/register", middlewares.RedirectIfAuthenticated, fc.ShowRegisterPage)

	r.GET("/logout", func(c *gin.Context) {
		c.SetCookie("Authorization", "", -1, "", "", false, true)

		c.HTML(http.StatusOK, "login.html", gin.H{
			"message": "Logged out",
		})
	})

	r.GET("/courses", middlewares.FERequireAuth, fc.GetCoursesPage)

	r.GET("/my-courses", middlewares.FERequireAuth, fc.GetMyCoursesPage)

	r.GET("/course/:id", middlewares.FERequireAuth, fc.GetCourseDetailPage)

	r.POST("/course/:id/purchase", middlewares.FERequireAuth, fc.BuyCourseFE)

	r.GET("/course/:id/modules", middlewares.FERequireAuth, fc.GetCourseModulesPage)
	r.GET("/course/:id/modules/:moduleId", middlewares.FERequireAuth, fc.GetCourseModulesPage)
	r.POST("/course/:id/modules/:moduleId/completion", middlewares.FERequireAuth, fc.ToggleModuleCompletion)

	r.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusNotFound, "404.html", gin.H{
			"title": "Page Not Found",
		})
	})
}

func moduleURL(courseID, moduleID uint) string {
	return fmt.Sprintf("/course/%d/modules/%d", courseID, moduleID)
}

func moduleURLWithType(courseID, moduleID uint, contentType string) string {
	return fmt.Sprintf("/course/%d/modules/%d?type=%s", courseID, moduleID, contentType)
}

func moduleCompletionURL(courseID, moduleID uint) string {
	return fmt.Sprintf("/course/%d/modules/%d/completion", courseID, moduleID)
}

func moduleTypeLabel(module models.ModuleWithIsCompleted) string {
	hasPDF := module.PDFContent != ""
	hasVideo := module.VideoContent != ""

	if hasPDF && hasVideo {
		return "PDF + Video"
	} else if hasPDF {
		return "PDF Content"
	} else if hasVideo {
		return "Video Content"
	}
	return "No Content"
}

func hasMultipleContent(module models.ModuleWithIsCompleted) bool {
	return module.PDFContent != "" && module.VideoContent != ""
}

func shouldShowPDF(module models.ModuleWithIsCompleted, contentType string) bool {
	if module.PDFContent == "" {
		return false
	}

	if hasMultipleContent(module) {
		return contentType != "video"
	}

	return true
}

func shouldShowVideo(module models.ModuleWithIsCompleted, contentType string) bool {
	if module.VideoContent == "" {
		return false
	}

	if hasMultipleContent(module) {
		return contentType == "video"
	}

	return true
}
