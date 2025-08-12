package handlers

import (
	"asset-service/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type FolderHandler struct {
	svc services.FolderService
}

func NewFolderHandler(svc services.FolderService) *FolderHandler {
	return &FolderHandler{svc: svc}
}

func (h *FolderHandler) CreateFolder(c *gin.Context) {
	var req struct {
		Name string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.svc.CreateFolder(req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, result)
}

func (h *FolderHandler) ListFolders(c *gin.Context) {
	result, err := h.svc.ListFolders()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *FolderHandler) GetFolderByID(c *gin.Context) {
	id := c.Param("id")

	result, err := h.svc.GetFolderByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *FolderHandler) DeleteFolder(c *gin.Context) {
	id := c.Param("id")

	err := h.svc.DeleteFolder(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
