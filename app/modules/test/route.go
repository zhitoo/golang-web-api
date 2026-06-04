package test

import (
	"github.com/gin-gonic/gin"
)

func Routes(r *gin.RouterGroup) {
	handler := NewHandler()
	r.GET("", handler.List)
	r.GET("/:id", handler.Show)
	r.POST("", handler.Store)
	r.POST("/header-binder-1", handler.HeaderBinder1)
	r.POST("/header-binder-2", handler.HeaderBinder2)
	r.POST("/query-binder-1", handler.QueryBinder1)
	r.POST("/query-binder-2", handler.QueryBinder2)
	r.POST("/uri-binder/:id/:name", handler.UriBinder)
	r.POST("/body-binder", handler.BodyBinder)
	r.POST("/file-binder", handler.FileBinder)
}
