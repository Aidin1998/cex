package main

import (
	"context"
	"errors"
	"os"
	"os/signal"

	"cex/cmd"
	"cex/internal/users"
	"cex/pkg/cfg"
	"cex/pkg/otel"
)

var (
	config    = cfg.MustLoad[Config]()
	zapLogger = cmd.NewZapLogger(config.IsDev)
	logger    = cmd.NewLogger(zapLogger)
	db        = cmd.CRDB(config.CockroachDB.DSN, logger.Handler())
)

func main() {
	defer zapLogger.Sync()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// Setup otel
	otelShutdown, err := otel.Setup(ctx)
	if err != nil {
		panic(err)
	}
	defer func() {
		err = errors.Join(err, otelShutdown(context.Background()))
	}()

	app := users.New(users.Opts{
		ListenAddress: "localhost:3000",
		DB:            db,
		Log:           logger,
	})

	app.Run()
}
