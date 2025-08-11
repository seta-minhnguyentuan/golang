package httpserver

import (
	"asset-service/internal/handlers"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RouterDeps struct {
	FolderService FolderService
}

type FolderService interface {
	Create(name string) (any, error)
	GetByID(id string) (any, error)
	List() ([]any, error)
	Delete(id string) error
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
		folders.POST("", h.Create)
		folders.GET("", h.List)
		folders.GET("/:id", h.GetByID)
		folders.DELETE("/:id", h.Delete)
	}

	return r
}
