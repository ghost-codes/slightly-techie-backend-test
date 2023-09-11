package config

import (
	"time"
)

type Config struct {
	AccessTokenDuration  time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
	SecretKey            string        `mapstructure:"SECRET_KEY"`
	DBDriver             string        `mapstructure:"DB_DRIVER"`
	ServerAddr           string        `mapstructure:"SEVER_ADDR"`
	ServerUrl            string        `mapstructure:"SEVER_URL"`
	SCHEME               string        `mapstructure:"SCHEME"`
	DBName               string        `mapstructure:"DB_NAME"`
}

func (config *Config) DBSource() string {
	return config.DBName
}
