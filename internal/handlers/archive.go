package handlers

import (
	"dbc/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ArchiveHandler struct {
	Service services.ArchiveService
}

func NewArchiveHandler(service services.ArchiveService) *ArchiveHandler {
	return &ArchiveHandler{Service: service}
}

func (h *ArchiveHandler) GetArchiveInformation(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}
	defer file.Close()

	info, err := h.Service.GetArchiveInfo(file, header.Filename)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, info)
}

func (h *ArchiveHandler) CreateArchive(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid form data"})
		return
	}

	archive, err := h.Service.CreateArchive(form.File["files[]"])
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Type", "application/zip")
	c.Header("Content-Disposition", `attachment; filename="archive.zip"`)
	c.Data(http.StatusOK, "application/zip", archive)
}
