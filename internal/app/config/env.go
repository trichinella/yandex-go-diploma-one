package config

import (
	"os"
)

func (config *Config) updateByEnv() {
	envGopherMartAddress := os.Getenv("RUN_ADDRESS")
	if envGopherMartAddress != "" {
		config.GopherMartAddress = envGopherMartAddress
	}

	envAccrualAddress := os.Getenv("ACCRUAL_SYSTEM_ADDRESS")
	if envAccrualAddress != "" {
		config.AccrualAddress = envAccrualAddress
	}

	envDatabaseDSN := os.Getenv("DATABASE_URI")
	if envDatabaseDSN != "" {
		config.DatabaseDSN = envDatabaseDSN
	}

	envJWTKey := os.Getenv("JWT_KEY")
	if envJWTKey != "" {
		config.JWTKey = envJWTKey
	}
}
