package httpserver

import (
	"context"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"

	"github.com/the-witcher-knight/jwt-encryption-server/internal/tracing"
)

func rootMiddleware(rootCtx context.Context) gin.HandlerFunc {
	return func(c *gin.Context) {
		tracer := tracing.FromContext(rootCtx).WithAttributes(
			tracing.String("host.name", c.Request.Host),
			tracing.String("url.path", c.Request.URL.String()),
			tracing.String("url.query", c.Request.URL.RawQuery),
			tracing.String("http.request.method", c.Request.Method),
			tracing.Int("http.request.body.size", int(c.Request.ContentLength)),
			tracing.String("http.request.proto", c.Request.Proto),
			tracing.String("http.request.remote_address", c.Request.RemoteAddr),
			tracing.String("user_agent.original", c.Request.UserAgent()),
		)

		defer func() {
			if p := recover(); p != nil {
				err, ok := p.(error)
				if !ok {
					err = fmt.Errorf("%v", p)
				}

				tracer.Error(err, "caught a panic", debug.Stack())
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error":             "internal server error",
					"error_description": "internal server error",
				})
			}
		}()

		updateReqCtx := func(ctx context.Context) context.Context {
			return tracing.SetInContext(ctx, tracer)
		}
		c.Request.WithContext(updateReqCtx(c.Request.Context()))

		// Go next step
		c.Next()

		tracer.WithAttributes(
			tracing.Int("http.response.status_code", c.Writer.Status()),
			tracing.Int("http.response.body.size", c.Writer.Size()),
		).Info("request handled successfully")
	}
}
