package routes

import (
	"dbc/internal/handlers"
	"dbc/internal/services"

	"github.com/gin-gonic/gin"
)

func InitRoutes(r *gin.Engine) {
	api := r.Group("/api")

	archiveService := services.NewArchiveService()
	emailService := services.NewEmailService()

	archiveHandler := handlers.NewArchiveHandler(archiveService)
	emailHandler := handlers.NewEmailHandler(emailService)
	{
		archive := api.Group("/archive")
		{
			archive.POST("/information", archiveHandler.GetArchiveInformation)
			archive.POST("/files", archiveHandler.CreateArchive)
		}
		api.POST("/mail/file", emailHandler.SendFileViaEmail)
	}
}
