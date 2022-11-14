package http

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/sync/errgroup"
)

const defaultPort = 8080

func TestConnectAndStartServiceWithContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	eg, ctx := errgroup.WithContext(ctx)

	// create server
	ConnectAndStartServiceWithContext(ctx, eg, defaultPort, simpleHandler())

	// wait to start listener
	time.Sleep(time.Second * 3)

	// check server handlers
	resp, err := sendRequest(fmt.Sprintf("http://0.0.0.0:8080%s", "/ping"))
	assert.Nil(t, err, "request ping err nil")
	assert.Equal(t, "pong", resp, "request ping")

	// close the ctx to shutdown the server
	cancel()

	if err := eg.Wait(); err != nil {
		t.Error("unexpected error")
	}
}

func sendRequest(url string) (answer string, err error) {
	resp, err := http.Get(url) // nolint: gosec
	if err != nil {
		return "", err
	}
	defer resp.Body.Close() // nolint

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), err
}

func simpleHandler() http.Handler {
	handler := http.NewServeMux()
	handler.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "pong") // nolint
	})

	return handler
}
