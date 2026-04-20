package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	CORS     CORSConfig
}

type ServerConfig struct {
	Port string
	Env  string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

type JWTConfig struct {
	Secret     string
	Expiration time.Duration
}

type CORSConfig struct {
	AllowOrigins []string
}

func Load() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	// .envファイルが存在しない場合は環境変数のみを使用
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	config := &Config{
		Server: ServerConfig{
			Port: getEnvOrDefault("SERVER_PORT", "8080"),
			Env:  getEnvOrDefault("ENV", "development"),
		},
		Database: DatabaseConfig{
			Host:     getEnvOrDefault("DB_HOST", "localhost"),
			Port:     getEnvOrDefault("DB_PORT", "5432"),
			User:     getEnvOrDefault("DB_USER", "attendance_user"),
			Password: getEnvOrDefault("DB_PASSWORD", "attendance_password"),
			Name:     getEnvOrDefault("DB_NAME", "attendance_db"),
			SSLMode:  getEnvOrDefault("DB_SSLMODE", "disable"),
		},
		JWT: JWTConfig{
			Secret:     getEnvOrDefault("JWT_SECRET", "your-secret-key"),
			Expiration: time.Duration(viper.GetInt("JWT_EXPIRATION_HOURS")) * time.Hour,
		},
		CORS: CORSConfig{
			AllowOrigins: viper.GetStringSlice("CORS_ALLOW_ORIGINS"),
		},
	}

	// JWT expiration のデフォルト設定
	if config.JWT.Expiration == 0 {
		config.JWT.Expiration = 24 * time.Hour
	}

	// CORS origins のデフォルト設定
	if len(config.CORS.AllowOrigins) == 0 {
		config.CORS.AllowOrigins = []string{"http://localhost:3000"}
	}

	return config, nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := viper.GetString(key); value != "" {
		return value
	}
	return defaultValue
}

func (c *DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.Name, c.SSLMode,
	)
}
