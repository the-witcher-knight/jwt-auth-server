package handler

import (
	"net/http"

	"github.com/the-witcher-knight/jwt-encryption-server/internal/httpserver"
)

var (
	errInvalidClientIDOrSecret = &httpserver.HTTPError{Code: http.StatusBadRequest, Message: "invalid client id or secret"}
)
