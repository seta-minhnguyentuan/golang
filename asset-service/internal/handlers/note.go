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

func (h *NoteHandler) ListNotes(c *gin.Context) {
	notes, err := h.NoteService.ListNotes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, notes)
}

func (h *NoteHandler) GetNote(c *gin.Context) {
	id := c.Param("id")
	note, err := h.NoteService.GetNote(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, note)
}

func (h *NoteHandler) UpdateNote(c *gin.Context) {
	var req struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := c.Param("id")
	updatedNote, err := h.NoteService.UpdateNote(id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedNote)
}

func (h *NoteHandler) DeleteNote(c *gin.Context) {
	id := c.Param("id")
	if err := h.NoteService.DeleteNote(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
