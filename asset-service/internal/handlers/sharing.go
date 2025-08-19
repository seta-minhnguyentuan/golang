package handlers

import (
	"asset-service/internal/models"
	"asset-service/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SharingHandler struct {
	svc services.SharingService
}

func NewSharingHandler(svc services.SharingService) *SharingHandler {
	return &SharingHandler{svc: svc}
}

type ShareRequest struct {
	UserID     uuid.UUID         `json:"userId" binding:"required"`
	Permission models.Permission `json:"permission" binding:"required"`
}

// Folder sharing handlers
func (h *SharingHandler) ShareFolder(c *gin.Context) {
	folderID, err := uuid.Parse(c.Param("folderId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid folder ID"})
		return
	}

	var req ShareRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user ID from JWT token context (set by auth middleware)
	ownerID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	ownerUUID, ok := ownerID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
		return
	}

	err = h.svc.ShareFolder(folderID, req.UserID, req.Permission, ownerUUID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Folder shared successfully"})
}

func (h *SharingHandler) RevokeFolderSharing(c *gin.Context) {
	folderID, err := uuid.Parse(c.Param("folderId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid folder ID"})
		return
	}

	userID, err := uuid.Parse(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Get owner ID from JWT token context
	ownerID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	ownerUUID, ok := ownerID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
		return
	}

	err = h.svc.RevokeFolderSharing(folderID, userID, ownerUUID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Folder sharing revoked successfully"})
}

func (h *SharingHandler) ListFolderSharings(c *gin.Context) {
	folderID, err := uuid.Parse(c.Param("folderId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid folder ID"})
		return
	}

	// Get owner ID from JWT token context
	ownerID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	ownerUUID, ok := ownerID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
		return
	}

	sharings, err := h.svc.ListFolderSharings(folderID, ownerUUID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, sharings)
}

// Note sharing handlers
func (h *SharingHandler) ShareNote(c *gin.Context) {
	noteID, err := uuid.Parse(c.Param("noteId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid note ID"})
		return
	}

	var req ShareRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user ID from JWT token context (set by auth middleware)
	ownerID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	ownerUUID, ok := ownerID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
		return
	}

	err = h.svc.ShareNote(noteID, req.UserID, req.Permission, ownerUUID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Note shared successfully"})
}

func (h *SharingHandler) RevokeNoteSharing(c *gin.Context) {
	noteID, err := uuid.Parse(c.Param("noteId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid note ID"})
		return
	}

	userID, err := uuid.Parse(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Get owner ID from JWT token context
	ownerID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	ownerUUID, ok := ownerID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
		return
	}

	err = h.svc.RevokeNoteSharing(noteID, userID, ownerUUID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Note sharing revoked successfully"})
}

func (h *SharingHandler) ListNoteSharings(c *gin.Context) {
	noteID, err := uuid.Parse(c.Param("noteId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid note ID"})
		return
	}

	// Get owner ID from JWT token context
	ownerID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	ownerUUID, ok := ownerID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
		return
	}

	sharings, err := h.svc.ListNoteSharings(noteID, ownerUUID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, sharings)
}
