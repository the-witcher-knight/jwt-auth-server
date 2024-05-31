package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/the-witcher-knight/jwt-encryption-server/internal/handler"
	"github.com/the-witcher-knight/jwt-encryption-server/internal/httpio"
	"github.com/the-witcher-knight/jwt-encryption-server/internal/logging"
	"github.com/the-witcher-knight/jwt-encryption-server/internal/service"
)

const (
	addr = ":8080"
)

func main() {
	ctx := context.Background()
	if err := run(ctx); err != nil {
		log.Printf("error %+v", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	logger, err := logging.NewLogger()
	if err != nil {
		return err
	}

	signatureService, err := service.NewRSASignatureService("private-key.pem")
	if err != nil {
		return err
	}

	hdl := handler.New(signatureService)
	mux := http.NewServeMux()

	addRoutes(mux, logger, hdl)

	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	return startServer(ctx, logger, mux)
}

func startServer(ctx context.Context, logger *logging.Logger, mux *http.ServeMux) error {
	srv := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	gctx, gcancel := context.WithCancelCause(ctx)
	wg := new(sync.WaitGroup)

	wg.Add(1)
	go func() {
		defer wg.Done()

		logger.Info(context.Background(), "starting server at port "+addr)
		if err := srv.ListenAndServe(); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				gcancel(err)
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		<-gctx.Done()

		logger.Info(context.Background(), "shutting down server")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			gcancel(err)
		}
	}()

	wg.Wait()
	if err := context.Cause(gctx); err != nil && !errors.Is(err, context.Canceled) {
		return err
	}

	return nil
}

func addRoutes(mux *http.ServeMux, logger *logging.Logger, hdl handler.Handler) {
	mux.Handle("/token", httpio.RootMiddleware(logger)(hdl.GenerateToken()))
	mux.Handle("/.well-known/jwks.json", httpio.RootMiddleware(logger)(hdl.GetJWKs()))
}
