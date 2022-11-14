package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/sync/errgroup"
)

// ConnectAndStartServiceWithContext return pointer to http.Service with context and WaitGroup
// have graceful shutdown of server in case of context.Done
func ConnectAndStartServiceWithContext(ctx context.Context, eg *errgroup.Group, port int, handler http.Handler) {
	srv := connectAndStartService(eg, port, handler)

	go func() {
		<-ctx.Done()

		cctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		srv.Shutdown(cctx) // nolint
	}()
}

// connectService simple create http.Service
// with address: "localhost:[port]" and provided Handler
func connectAndStartService(eg *errgroup.Group, port int, handler http.Handler) *http.Server {
	srv := &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%d", port),
		Handler: handler,
	}

	eg.Go(func() error {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			return err
		}
		return nil
	})

	return srv
}
