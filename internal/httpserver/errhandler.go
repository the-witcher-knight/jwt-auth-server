package httpserver

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/the-witcher-knight/jwt-encryption-server/internal/tracing"
)

func ErrorHandler(fn func(*gin.Context) error) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := fn(c); err != nil {
			var httpErr *HTTPError
			if errors.As(err, &httpErr) {
				c.JSON(httpErr.Code, httpErr)
				return
			}

			tracing.FromContext(c).Error(err, "internal server error")
			c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]string{
				"error":             "internal server error",
				"error_description": "internal server error",
			})
		}
	}
}
