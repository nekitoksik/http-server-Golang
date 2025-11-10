package config

import (
	"errors"
	"log"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
}

type ServerConfig struct {
	Address  string
	LogLevel string
}

type DatabaseConfig struct {
	URL            string
	MigrationsPath string
}

type JWTConfig struct {
	Secret               string
	AccessTokenDuration  time.Duration
	RefreshTokenDuration time.Duration
}

func LoadConfig() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	setDefaults()

	if err := viper.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			log.Println("No .env file found, using defaults settings")
		} else {
			log.Printf("Error reading config file: %v", err)
		}
	}

	cfg := &Config{
		Server: ServerConfig{
			Address:  viper.GetString("SERVER_ADDRESS"),
			LogLevel: viper.GetString("LOG_LEVEL"),
		},
		Database: DatabaseConfig{
			URL:            viper.GetString("DATABASE_URL"),
			MigrationsPath: viper.GetString("MIGRATIONS_PATH"),
		},
		JWT: JWTConfig{
			Secret:               viper.GetString("JWT_SECRET"),
			RefreshTokenDuration: viper.GetDuration("ACCESS_TOKEN_CODE_EXPIRY"),
			AccessTokenDuration:  viper.GetDuration("REFRESH_TOKEN_DURATION"),
		},
	}

	if err := validateConfig(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func setDefaults() {
	viper.SetDefault("SERVER_ADDRES", ":8080")
	viper.SetDefault("LOG_LEVEL", "info")
	viper.SetDefault("MIGRATIONS_PATH", "internal/db/migrations")
}

func validateConfig(cfg *Config) error {
	if cfg.Database.URL == "" {
		return errors.New("DATABASE_URL is required field")
	}

	if cfg.Database.MigrationsPath == "" {
		return errors.New("MIGRATIONS_PATH is required field")
	}

	return nil
}
