package config

import (
	"fmt"
	"os"
)

/*
 * アプリケーション設定
 * 環境変数から設定を読み込む
 */

// Config アプリケーション設定
type Config struct {
	DatabaseURL string
	ServerPort  string
	JWTSecret   string
}

// LoadConfig 環境変数から設定を読み込む
func LoadConfig() (*Config, error) {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return nil, fmt.Errorf("DATABASE_URL が設定されていません")
	}

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET が設定されていません")
	}

	return &Config{
		DatabaseURL: dbURL,
		ServerPort:  port,
		JWTSecret:   jwtSecret,
	}, nil
}
