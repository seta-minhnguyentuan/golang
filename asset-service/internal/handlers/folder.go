package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type FolderService interface {
	Create(name string) (any, error)
	GetByID(id string) (any, error)
	List() ([]any, error)
	Delete(id string) error
}

type FolderHandler struct {
	svc FolderService
}

func NewFolderHandler(svc FolderService) *FolderHandler {
	return &FolderHandler{svc: svc}
}

func (h *FolderHandler) Create(c *gin.Context) {
	var req struct {
		Name string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.svc.Create(req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, result)
}

func (h *FolderHandler) List(c *gin.Context) {
	result, err := h.svc.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *FolderHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	result, err := h.svc.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *FolderHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	err := h.svc.Delete(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
