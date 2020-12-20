package config

import (
	"os"

	"go.uber.org/zap"
)

// CadenceConfig ...
type CadenceConfig struct {
	Domain   string
	Service  string
	HostPort string
}

//AppConfig ...
type AppConfig struct {
	Env            string
	WorkerTaskList string
	Cadence        CadenceConfig
	Logger         *zap.Logger
}

// LoadConfig setup the config for the code run
func (h *AppConfig) LoadConfig(configPath string) {
	h.Cadence.Domain = os.Getenv("CADENCE_DOMAIN")
	h.Cadence.HostPort = os.Getenv("CADENCE_HOST")
	h.Cadence.Service = os.Getenv("CADENCE_SERVICE")

	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	h.Logger = logger
	logger.Debug("Finished loading Configuration!")
}

