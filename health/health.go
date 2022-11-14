package health1

import (
	"context"
	"net/http"

	chassishttp "github.com/kurennoy16/chassis-example/http"
	"golang.org/x/sync/errgroup"
)

const PortDefault = 8888

func ConnectHealthCheckServiceWithContext(ctx context.Context,
	eg *errgroup.Group,
	port int,
	healthChecks []func() error,
) {
	chassishttp.ConnectAndStartServiceWithContext(ctx, eg, port, handler(healthChecks))
}

func handler(healthChecks []func() error) http.Handler {
	hdl := http.NewServeMux()
	hdl.HandleFunc("/", func(res http.ResponseWriter, _ *http.Request) {
		serveCheck(res, healthChecks)
	})

	return hdl
}

func serveCheck(response http.ResponseWriter, checks []func() error) {
	var writtenHeader bool
	for _, check := range checks {
		if err := check(); err != nil {
			if !writtenHeader {
				response.WriteHeader(http.StatusInternalServerError)
				writtenHeader = true
			}
			response.Write([]byte(err.Error())) // nolint
			response.Write([]byte("\n\n"))      // nolint
		}
	}

	if !writtenHeader {
		response.WriteHeader(http.StatusNoContent)
	}
}
