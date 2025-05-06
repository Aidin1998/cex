package users

import (
	"log/slog"

	"cex/internal/users/api"

	"gorm.io/gorm"
)

type Opts struct {
	Log           *slog.Logger
	DB            *gorm.DB
	ListenAddress string
}

type App struct {
	log           *slog.Logger
	db            *gorm.DB
	listenAddress string

	api *api.API
}

func New(opts Opts) *App {
	return &App{
		api:           api.New(opts.Log),
		log:           opts.Log,
		db:            opts.DB,
		listenAddress: opts.ListenAddress,
	}
}

func (a *App) Run() {
	a.api.Serve(a.listenAddress)
}
