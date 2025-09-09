package handlers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB) {
	api := r.Group("/api")
	{
		api.POST("/upload/resume", func(c *gin.Context) { UploadResume(c, db) })
		api.POST("/upload/vacancy", func(c *gin.Context) { UploadVacancy(c, db) })
		api.POST("/analyze", func(c *gin.Context) { AnalyzeResume(c, db) })
		api.GET("/health", HealthCheck)

		api.GET("/vacancies", func(c *gin.Context) { GetVacancies(c, db) })
		api.POST("/analyze-resume", func(c *gin.Context) { AnalyzeResumeForVacancy(c, db) })
	}

	r.Static("/static", "./frontend")
	r.GET("/", func(c *gin.Context) {
		c.File("./frontend/index.html")
	})
	r.GET("/interview.html", func(c *gin.Context) {
		c.File("./frontend/interview.html")
	})

	r.GET("/health", HealthCheck)
}

func SetupResumeRoutes(r *gin.Engine, db *gorm.DB) {
	r.POST("/upload/resume", func(c *gin.Context) { UploadResume(c, db) })
	r.POST("/upload/vacancy", func(c *gin.Context) { UploadVacancy(c, db) })
	r.POST("/analyze", func(c *gin.Context) { AnalyzeResume(c, db) })
	r.GET("/health", HealthCheck)

	r.GET("/api/vacancies", func(c *gin.Context) { GetVacancies(c, db) })
	r.POST("/api/analyze-resume", func(c *gin.Context) { AnalyzeResumeForVacancy(c, db) })
}

func HealthCheck(c *gin.Context) {
	c.JSON(200, gin.H{"status": "ok"})
}
