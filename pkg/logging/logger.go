package logging

import "github.com/zhitoo/golang-web-api/config"

// ScopedLogger is a logger with category and subcategory pre-set.
// Obtain one via Logger.With(cat, sub) to avoid repeating the same
// category pair on every call inside a package or service.
type ScopedLogger interface {
	Debug(string, map[ExtraKey]any)
	Debugf(string, ...any)
	Info(string, map[ExtraKey]any)
	Infof(string, ...any)
	Warn(string, map[ExtraKey]any)
	Warnf(string, ...any)
	Error(string, map[ExtraKey]any)
	Errorf(string, ...any)
	Fatal(string, map[ExtraKey]any)
	Fatalf(string, ...any)
}

type Logger interface {
	Init()
	With(Category, SubCategory) ScopedLogger

	Debug(Category, SubCategory, string, map[ExtraKey]any)
	Debugf(string, ...any)

	Info(Category, SubCategory, string, map[ExtraKey]any)
	Infof(string, ...any)

	Warn(Category, SubCategory, string, map[ExtraKey]any)
	Warnf(string, ...any)

	Error(Category, SubCategory, string, map[ExtraKey]any)
	Errorf(string, ...any)

	Fatal(Category, SubCategory, string, map[ExtraKey]any)
	Fatalf(string, ...any)
}

func NewLogger(cfg *config.Config) Logger {
	return newZeroLog(cfg)
}
