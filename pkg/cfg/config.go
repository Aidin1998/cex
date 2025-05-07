package cfg

import (
	"os"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

// AccountsConfig holds HTTP port & DSN for accounts module
type AccountsConfig struct {
	Port string `mapstructure:"ACCOUNTS_PORT" validate:"required"`
	DSN  string `mapstructure:"ACCOUNTS_DSN"  validate:"required,url"`
}

// DBConfig holds the global database connection string.
type DBConfig struct {
	// URL is the Postgres (CockroachDB) DSN.
	URL string `mapstructure:"url"`
}

// HTTPConfig holds default HTTP server settings.
type HTTPConfig struct {
	// Port is the TCP port for the HTTP gateway.
	Port string `mapstructure:"port"`
}

// UsersConfig holds settings for the users module.
type UsersConfig struct {
	// DSN is the Postgres DSN for the users service.
	DSN string `mapstructure:"dsn"`
	// Port is the HTTP port for the users service.
	Port string `mapstructure:"port"`
	// JWTSecret is the HMAC secret for signing user JWTs.
	JWTSecret string `mapstructure:"jwt_secret"`
}

type Config struct {
	Accounts AccountsConfig
	Queue    struct {
		Topics []string
		URL    string
	}
	Kafka struct {
		Brokers       []string
		TopicAccounts string
	}
	DB    DBConfig    `mapstructure:"db"`
	HTTP  HTTPConfig  `mapstructure:"http"`
	Users UsersConfig `mapstructure:"users"`
}

// Cfg is the singleton instance loaded by Init.
var Cfg Config

// Init reads config.(yaml|json|toml|env) into Cfg
func Init() {
	v := viper.New()
	v.AutomaticEnv()

	if err := v.Unmarshal(&Cfg); err != nil {
		panic("failed to unmarshal config: " + err.Error())
	}

	if err := validator.New().Struct(Cfg.Accounts); err != nil {
		panic("invalid accounts config: " + err.Error())
	}

	// ─── Fail‐fast checks ───────────────────────────────
	if Cfg.DB.URL == "" {
		panic("config: db.url is required")
	}
	// ────────────────────────────────────────────────────
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
