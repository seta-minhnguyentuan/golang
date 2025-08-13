package handlers

import (
	"asset-service/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type NoteHandler struct {
	NoteService services.NoteService
}

func NewNoteHandler(noteService services.NoteService) *NoteHandler {
	return &NoteHandler{
		NoteService: noteService,
	}
}

func (h *NoteHandler) CreateNote(c *gin.Context) {
	var req struct {
		Title    string    `json:"title"`
		Content  string    `json:"content"`
		FolderID uuid.UUID `json:"folder_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	note, err := h.NoteService.CreateNote(req.Title, req.Content, req.FolderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, note)
}
