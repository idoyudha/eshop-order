package config

import "github.com/ilyakaznacheev/cleanenv"

type (
	Config struct {
		App  `yaml:"app"`
		HTTP `yaml:"http"`
		Log  `yaml:"log"`
		PostgreSQLCommand
		PostgreSQLQuery
		AuthService
		Kafka
	}

	App struct {
		Name    string `env-required:"true" yaml:"name" env:"APP_NAME"`
		Version string `env-required:"true" yaml:"version" env:"APP_VERSION"`
	}

	HTTP struct {
		Port string `env-required:"true" yaml:"port" env:"HTTP_PORT"`
	}

	PostgreSQLCommand struct {
		URL         string `env-required:"true" env:"POSTGRESQL_COMMAND_URL"`
		ConnTimeout int    `env-required:"true" env:"POSTGRESQL_COMMAND_CONN_TIMEOUT"`
		ConnAttemps int    `env-required:"true" env:"POSTGRESQL_COMMAND_CONN_ATTEMPS"`
		MaxPoolSize int    `env-required:"true" env:"POSTGRESQL_COMMAND_MAX_POOL_SIZE"`
	}

	PostgreSQLQuery struct {
		URL         string `env-required:"true" env:"POSTGRESQL_QUERY_URL"`
		ConnTimeout int    `env-required:"true" env:"POSTGRESQL_QUERY_CONN_TIMEOUT"`
		ConnAttemps int    `env-required:"true" env:"POSTGRESQL_QUERY_CONN_ATTEMPS"`
		MaxPoolSize int    `env-required:"true" env:"POSTGRESQL_QUERY_MAX_POOL_SIZE"`
	}

	Log struct {
		Level string `env-required:"true" yaml:"level" env:"LOG_LEVEL"`
	}

	AuthService struct {
		BaseURL string `env-required:"true" env:"AUTH_SERVICE"`
	}

	Kafka struct {
		Broker string `env-required:"true" env:"KAFKA_BROKER"`
	}
)

func NewConfig() (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadConfig("./config/config.yml", cfg)
	if err != nil {
		return nil, err
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
