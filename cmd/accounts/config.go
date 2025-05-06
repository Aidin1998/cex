package accounts

import (
	"cex/pkg/cfg"
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

var Cfg = cfg.Cfg

type Config struct {
	IsDev bool

	CockroachDB struct {
		DSN string
	}
}

func Init() {
	v := viper.New()
	v.AutomaticEnv()
	v.SetDefault("ACCOUNTS_PORT", "8081")
	if err := v.Unmarshal(&Cfg); err != nil {
		log.Fatal("failed to unmarshal config", "error", err)
	}
	if err := validator.New().Struct(Cfg.Accounts); err != nil {
		log.Fatal("invalid accounts config", "error", err)
	}
}
