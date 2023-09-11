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
	DBHOST               string        `mapstructure:"DB_HOST"`
	DBPort               string        `mapstructure:"DB_PORT"`
	DBUser               string        `mapstructure:"DB_USER"`
	DBPassword           string        `mapstructure:"DB_PASSWORD"`
	DBName               string        `mapstructure:"DB_NAME"`
	PaystackSK           string        `mapstructure:"PAYSTACK_SECRET_KEY"`
	PaystackPK           string        `mapstructure:"PAYSTACK_PUBLIC_KEY"`
}

func (config *Config) DBSource() string {
	return config.DBName
}
