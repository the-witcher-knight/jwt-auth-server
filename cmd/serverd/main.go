package main

import (
	"context"
	"crypto/rsa"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/the-witcher-knight/jwt-encryption-server/internal/handler"
	"github.com/the-witcher-knight/jwt-encryption-server/internal/httpserver"
	"github.com/the-witcher-knight/jwt-encryption-server/internal/secrets"
	"github.com/the-witcher-knight/jwt-encryption-server/internal/service"
	"github.com/the-witcher-knight/jwt-encryption-server/internal/tracing"
)

const (
	addr           = ":8080"
	privateKeyPath = "private-key.pem"
)

var (
	logger = log.New(os.Stdout, "", log.LstdFlags|log.Lmicroseconds)
)

func main() {
	if err := run(); err != nil {
		logger.Printf("server exited abnormally %+v", err)
		os.Exit(1)
	}
}

func run() error {
	fileBytes, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return err
	}

	privateKey, err := secrets.LoadPrivateKeyFromPEM[*rsa.PrivateKey](fileBytes, "")
	if err != nil {
		return err
	}

	tracer := tracing.New()
	defer func() {
		if err := tracer.Flush(); err != nil {
			logger.Printf("error flushing tracer: %+v", err)
		}
	}()

	svc, err := service.NewRSASignatureService(privateKey)
	hdl := handler.New(svc)

	// Setup HTTP server
	srv := &http.Server{
		Addr:         addr,
		ReadTimeout:  time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      newHTTPHandler(tracing.SetInContext(context.Background(), tracer), hdl),
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	srvErr := make(chan error, 1)
	go func() {
		srvErr <- srv.ListenAndServe()
	}()

	// Wait for interruption
	select {
	case err := <-srvErr:
		// Error when starting HTTP server
		return err
	case <-ctx.Done():
		stop()
	}

	return srv.Shutdown(context.Background())
}

func newHTTPHandler(rootCtx context.Context, hdl handler.Handler) http.Handler {
	router := httpserver.NewRouter(rootCtx)

	// Register handlers
	router.POST("/token", hdl.GenerateToken())
	router.GET("/.well-known/jwks.json", hdl.GetJWKs())

	return router
}
