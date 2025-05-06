package cfg

import (
	"os"
	"strings"

	"github.com/spf13/viper"
)

type accountsConfig struct {
	Port      string `mapstructure:"port"`
	JWTSecret string `mapstructure:"jwt_secret"`
}

type Config struct {
	DB struct {
		URL string
	}
}

var Cfg Config

// Init reads config.(yaml|json|toml|env) into Cfg
func Init() {
	viper.SetConfigName("config") // looks for config.(yaml|json|toml) in working dir
	viper.AddConfigPath(".")
	viper.AutomaticEnv() // also read from $ENV

	if err := viper.ReadInConfig(); err != nil {
		panic("failed to read config: " + err.Error())
	}
	if err := viper.Unmarshal(&Cfg); err != nil {
		panic("failed to unmarshal config: " + err.Error())
	}
}

func MustLoad[T any]() *T {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		panic("Couldn't load configuration, cannot start. Terminating. Error: " + err.Error())
	}

	for _, k := range viper.AllKeys() {
		value := viper.GetString(k)
		if strings.HasPrefix(value, "${") && strings.HasSuffix(value, "}") {
			viper.Set(k, getEnvOrPanic(strings.TrimSuffix(strings.TrimPrefix(value, "${"), "}")))
		}
	}

	var config T
	if err := viper.Unmarshal(&config); err != nil {
		panic("Failed to unmarshal config file: " + err.Error())
	}
	return &config
}

func getEnvOrPanic(env string) string {
	res := os.Getenv(env)
	if len(res) == 0 {
		panic("Mandatory env variable not found:" + env)
	}
	return res
}
