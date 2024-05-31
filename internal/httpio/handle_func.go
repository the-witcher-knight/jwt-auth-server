package httpio

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/the-witcher-knight/jwt-encryption-server/internal/logging"
)

func HandleFunc(fn func(http.ResponseWriter, *http.Request) error) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		w.Header().Set("Content-Type", "application/json")

		if err := fn(w, r); err != nil {
			var apiErr APIError
			if errors.As(err, &apiErr) {
				w.WriteHeader(apiErr.Status)
				_ = WriteJSON(w, Message{
					Message: apiErr.Error(),
				})

				return
			}

			w.WriteHeader(http.StatusInternalServerError)
			_ = WriteJSON(w, Message{
				Message: "internal server error",
			})

			logging.MayBeOf(ctx).
				Error(context.Background(), err, "Got error while processing request")
		}
	})
}

func ReadJSON[T any](reader io.Reader) (T, error) {
	var t T
	if err := json.NewDecoder(reader).Decode(&t); err != nil {
		return *new(T), err
	}

	return t, nil
}

func WriteJSON[T any](w io.Writer, t T) error {
	if err := json.NewEncoder(w).Encode(t); err != nil {
		return err
	}

	return nil
}
