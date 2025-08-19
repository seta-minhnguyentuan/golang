package httpserver

import (
	"asset-service/internal/handlers"
	"asset-service/internal/services"
	"net/http"
	"shared/middlewares"

	"github.com/gin-gonic/gin"
)

type RouterDeps struct {
	FolderService services.FolderService
	NoteService   services.NoteService
}

func NewRouter(deps RouterDeps) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	v1 := r.Group("/api/v1")

	folders := v1.Group("/folders")
	folders.Use(middlewares.AuthMiddleware())
	{
		h := handlers.NewFolderHandler(deps.FolderService)
		folders.POST("", h.CreateFolder)
		folders.GET("", h.ListFolders)
		folders.GET("/:id", h.GetFolderByID)
		folders.DELETE("/:id", h.DeleteFolder)
	}

	notes := v1.Group("/notes")
	{
		h := handlers.NewNoteHandler(deps.NoteService)
		notes.POST("", h.CreateNote)
		notes.GET("", h.ListNotes)
		notes.GET("/:id", h.GetNote)
		notes.PUT("/:id", h.UpdateNote)
		notes.DELETE("/:id", h.DeleteNote)
	}

	return r
}
