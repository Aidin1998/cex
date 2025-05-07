package metrics

import (
	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Middleware sets up Prometheus metrics middleware.
func Middleware() echo.MiddlewareFunc {
	return echoprometheus.NewMiddleware("namespace")
}

// RegisterMetricsEndpoint registers the /metrics endpoint.
func RegisterMetricsEndpoint(e *echo.Echo) {
	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))
}
