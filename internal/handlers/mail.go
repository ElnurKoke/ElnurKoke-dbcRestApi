package handlers

import (
	"net/http"
	"os"
	"strings"

	"dbc/internal/services"
	"dbc/internal/utils"

	"github.com/gin-gonic/gin"
)

type EmailHandler struct {
	Service services.EmailService
}

func NewEmailHandler(service services.EmailService) *EmailHandler {
	return &EmailHandler{Service: service}
}

func (h *EmailHandler) SendFileViaEmail(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}
	defer file.Close()

	emails := c.PostForm("emails")
	if emails == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Emails are required"})
		return
	}

	validTypes := []string{"application/pdf", "application/vnd.openxmlformats-officedocument.wordprocessingml.document"}
	mimeType := utils.GetMimeType(header.Filename)
	if !utils.IsValidMimeType(mimeType, validTypes) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type"})
		return
	}

	emailList := strings.Split(emails, ",")
	err = h.Service.SendEmail(os.Getenv("SMTP_USER"), os.Getenv("SMTP_PASSWORD"), emailList, header.Filename, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Email sent successfully"})
}
