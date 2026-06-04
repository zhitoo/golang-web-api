package logging

import (
	"carsale/config"
	"path/filepath"
	"runtime"
)

func getLogPath(cfg *config.Config) string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(filename), "..", "..", "logs", cfg.Logger.FileName)
}
