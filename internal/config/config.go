package config

import (
	"flag"
	"os"

	"github.com/joho/godotenv"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	ServerAddress string
	DBDSN         string
	LogLevel      string
}

func GetConfig() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{}

	flag.StringVar(&cfg.ServerAddress, "a", "localhost:8080", "address and port to run server")

	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		cfg.ServerAddress = envRunAddr
	}

	if envDBConn := os.Getenv("DBDSN"); envDBConn != "" {
		cfg.DBDSN = envDBConn
	} else {
		flag.StringVar(&cfg.DBDSN, "d", "", "DBDSN for database")
	}

	if envLogLevel := os.Getenv("LOG_LEVEL"); envLogLevel != "" {
		cfg.LogLevel = envLogLevel
	} else {
		cfg.LogLevel = zapcore.ErrorLevel.String()
	}

	flag.Parse()

	return cfg, nil
}
