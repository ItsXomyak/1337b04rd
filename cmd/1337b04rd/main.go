package main

import (
	"1337b04rd/config"
	"1337b04rd/internal/common/logger"
)

func main() {
	cfg := config.Load()
	logger.Init(cfg.AppEnv)

	logger.Info("App started", "env", cfg.AppEnv, "port", cfg.Port)
}
