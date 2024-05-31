package handler

import (
	"net/http"

	"github.com/the-witcher-knight/jwt-encryption-server/internal/httpio"
)

var (
	errInvalidClientIDOrSecret = httpio.APIError{Status: http.StatusBadRequest, Message: "invalid client id or secret"}
)
