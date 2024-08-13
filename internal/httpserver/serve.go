package httpserver

import (
	"context"
	"errors"

	"github.com/gin-gonic/gin"
)

func NewRouter(ctx context.Context) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	router.Use(rootMiddleware(ctx))
	blankGroup(router)

	return router
}

func blankGroup(router *gin.Engine) {
	blank := router.Group("/_")
	blank.GET("/panic", func(c *gin.Context) {
		panic(errors.New("simulate panic"))
	})
}
