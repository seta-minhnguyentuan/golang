package httpserver

import (
	"asset-service/internal/handlers"
	"asset-service/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RouterDeps struct {
	FolderService services.FolderService
}

func NewRouter(deps RouterDeps) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	v1 := r.Group("/api/v1")

	folders := v1.Group("/folders")
	{
		h := handlers.NewFolderHandler(deps.FolderService)
		folders.POST("", h.CreateFolder)
		folders.GET("", h.ListFolders)
		folders.GET("/:id", h.GetFolderByID)
		folders.DELETE("/:id", h.DeleteFolder)
	}

	return r
}
