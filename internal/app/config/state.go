package config

import (
	"flag"
	"sync"
)

var once sync.Once
var cfg *Config

func State() *Config {
	once.Do(func() {
		flag.Parse()

		cfg = &Config{}
		cfg.updateByFlags(baseOptions)
		cfg.updateByEnv()
	})

	return cfg
}
