package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/the-witcher-knight/jwt-encryption-server/internal/httpio"
	"github.com/the-witcher-knight/jwt-encryption-server/internal/service"
)

const (
	clientID     string = "sample-client-id"
	clientSecret string = "sample-client-secret"
)

var (
	payload = map[string]interface{}{
		"iss": "http://localhost:8080/",
		"sub": "sample_subject@clients",
		"aud": "http://localhost:9999/",
		"iat": 1717176856,
		"exp": 1717263256,
		"gty": "client-credentials",
		"azp": "sample_subject",
	}
)

type Handler struct {
	srv service.SignatureService
}

func New(srv service.SignatureService) Handler {
	return Handler{
		srv: srv,
	}
}

func (h Handler) GenerateToken() http.Handler {
	return httpio.HandleFunc(func(w http.ResponseWriter, r *http.Request) error {
		req, err := httpio.ReadJSON[generateTokenRequest](r.Body)
		if err != nil {
			return err
		}

		if req.ClientID != clientID || req.ClientSecret != clientSecret {
			return errInvalidClientIDOrSecret
		}

		payloadStr, err := json.Marshal(payload)
		if err != nil {
			return fmt.Errorf("could not marshal payload: %w", err)
		}

		token, err := h.srv.GenerateToken(payloadStr)
		if err != nil {
			return err
		}

		return httpio.WriteJSON(w, generateTokenResponse{
			AccessToken: token,
		})
	})
}

func (h Handler) GetJWKs() http.Handler {
	return httpio.HandleFunc(func(w http.ResponseWriter, r *http.Request) error {
		jwks, err := h.srv.GetJWKs()
		if err != nil {
			return err
		}

		return httpio.WriteJSON(w, jwks)
	})
}
