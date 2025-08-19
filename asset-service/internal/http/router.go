package httpserver

import (
	"asset-service/internal/handlers"
	"asset-service/internal/services"
	"net/http"
	"shared/middlewares"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type RouterDeps struct {
	FolderService  services.FolderService
	NoteService    services.NoteService
	SharingService services.SharingService
}

func NewRouter(deps RouterDeps) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())

	// CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:5173", "http://localhost:5174"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

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
		folders.GET("/:folderId", h.GetFolderByID)
		folders.DELETE("/:folderId", h.DeleteFolder)

		// Folder sharing endpoints
		sharingHandler := handlers.NewSharingHandler(deps.SharingService)
		folders.POST("/:folderId/share", sharingHandler.ShareFolder)
		folders.DELETE("/:folderId/share/:userId", sharingHandler.RevokeFolderSharing)
		folders.GET("/:folderId/share", sharingHandler.ListFolderSharings)
	}

	notes := v1.Group("/notes")
	notes.Use(middlewares.AuthMiddleware())
	{
		h := handlers.NewNoteHandler(deps.NoteService)
		notes.POST("", h.CreateNote)
		notes.GET("", h.ListNotes)
		notes.GET("/:noteId", h.GetNote)
		notes.PUT("/:noteId", h.UpdateNote)
		notes.DELETE("/:noteId", h.DeleteNote)

		// Note sharing endpoints
		sharingHandler := handlers.NewSharingHandler(deps.SharingService)
		notes.POST("/:noteId/share", sharingHandler.ShareNote)
		notes.DELETE("/:noteId/share/:userId", sharingHandler.RevokeNoteSharing)
		notes.GET("/:noteId/share", sharingHandler.ListNoteSharings)
	}

	return r
}
