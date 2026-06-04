package logging

import (
	"path/filepath"
	"runtime"

	"github.com/zhitoo/golang-web-api/config"
)

func getLogPath(cfg *config.Config) string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(filename), "..", "..", "logs", cfg.Logger.FileName)
}
