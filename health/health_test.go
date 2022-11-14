package health1

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"
)

func TestConnectHealthCheckServiceWithContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var checks []func() error
	eg, ctx := errgroup.WithContext(ctx)

	ConnectHealthCheckServiceWithContext(ctx, eg, PortDefault, checks)

	t.Run("check live server", func(t *testing.T) {
		resp, err := http.Get("http://localhost:8888/")
		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})

	t.Run("server shutdowned after ctx cancel", func(t *testing.T) {
		cancel()
		var closed = make(chan struct{})
		go func() {
			assert.NoError(t, eg.Wait())
			close(closed)
		}()
		select {
		case <-time.After(3 * time.Second):
			t.Fatal("server was not closed after 3s timeout")
		case <-closed:
			resp, err := http.Get("http://localhost:8888/")
			assert.Nil(t, resp)
			require.Error(t, err)
			assert.True(t, strings.Contains(err.Error(), "connect: connection refused"))
		}
	})
}

func Test_ServeCheck(t *testing.T) {
	tt := map[string]struct {
		checks  []func() error
		expCode int
	}{
		"error": {
			checks: []func() error{
				func() error {
					return fmt.Errorf("something went wrong")
				},
			},
			expCode: http.StatusInternalServerError,
		},
		"no_error": {
			checks: []func() error{
				func() error {
					return nil
				},
			},
			expCode: http.StatusNoContent,
		},
		"multiple statuses": {
			checks: []func() error{
				func() error {
					return nil
				},
				func() error {
					return fmt.Errorf("something went wrong")
				},
			},
			expCode: http.StatusInternalServerError,
		},
	}

	for testName, testCase := range tt {
		t.Run(testName, func(t *testing.T) {
			w := httptest.NewRecorder()
			serveCheck(w, testCase.checks)

			assert.Equal(t, w.Code, testCase.expCode)
		})
	}
}
