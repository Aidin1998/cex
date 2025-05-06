package api

import (
	"log/slog"

	"cex/pkg/apiutil"

	"github.com/labstack/echo/v4"
)

type API struct {
	log *slog.Logger
}

func New(log *slog.Logger) *API {
	return &API{
		log: log,
	}
}

func (a *API) Serve(listenAddr string) error {
	e := echo.New()

	e.HideBanner = true
	e.HidePort = true
	e.HTTPErrorHandler = apiutil.ErrorHandler(a.log)
	//e.Validator = &Validator{validator.New()}

	//e.GET("/docs/*", echoSwagger.WrapHandler)

	//v1 := e.Group("/api/v1")
	//a.registerPublicRoutes(v1)

	//authenticated := v1.Group("", a.auth.AuthMiddleware)
	//a.registerAuthenticatedRoutes(authenticated)

	a.log.Info("http server listening on " + listenAddr)
	return e.Start(listenAddr)
}
