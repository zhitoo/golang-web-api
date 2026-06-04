package logging

import "github.com/zhitoo/golang-web-api/config"

type Logger interface {
	Init()

	Debug(Categoty, SubCategoty, string, map[ExtraKey]any)
	Debugf(string, ...any)

	Info(Categoty, SubCategoty, string, map[ExtraKey]any)
	Infof(string, ...any)

	Warn(Categoty, SubCategoty, string, map[ExtraKey]any)
	Warnf(string, ...any)

	Error(Categoty, SubCategoty, string, map[ExtraKey]any)
	Errorf(string, ...any)

	Fatal(Categoty, SubCategoty, string, map[ExtraKey]any)
	Fatalf(string, ...any)
}

func NewLogger(cfg *config.Config) Logger {
	return newZeroLog(cfg)
}
