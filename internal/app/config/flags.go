package config

import (
	"flag"
)

type options struct {
	GopherMartAddress string
	AccrualAddress    string
	DatabaseDSN       string
	JWTKey            string
}

var baseOptions = options{}

func init() {
	flag.StringVar(&baseOptions.GopherMartAddress, "a", "", "GopherMart host")
	flag.StringVar(&baseOptions.DatabaseDSN, "d", "", "DSN for database")
	flag.StringVar(&baseOptions.AccrualAddress, "r", "", "Accrual host")
	flag.StringVar(&baseOptions.JWTKey, "jk", "simple_test_secret_key", "JWT key")
}

func (config *Config) updateByFlags(o options) {
	config.GopherMartAddress = o.GopherMartAddress
	config.AccrualAddress = o.AccrualAddress
	config.JWTKey = o.JWTKey
	config.DatabaseDSN = o.DatabaseDSN
}
