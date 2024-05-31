package httpio

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/the-witcher-knight/jwt-encryption-server/internal/logging"
)

func RootMiddleware(logger *logging.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx = logging.SetInCtx(ctx, logger)

			defer func() {
				if p := recover(); p != nil {
					err, ok := p.(error)
					if !ok {
						err = fmt.Errorf("%+v", p)
					}

					logger.Error(ctx, err, "caught panic",
						logging.AttributeString("stacktrace", string(debug.Stack())),
					)
				}
			}()

			next.ServeHTTP(w, r.WithContext(ctx))

			logger.With(
				logging.AttributeString("method", r.Method),
				logging.AttributeString("path", r.URL.Path),
				logging.AttributeString("host", r.Host),
			).Info(ctx, "Served request")
		})
	}
}
