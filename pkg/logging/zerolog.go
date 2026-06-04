package logging

import (
	"carsale/config"
	"log"
	"os"
	"sync"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

var zeroLogLevelMap = map[string]zerolog.Level{
	"debug": zerolog.DebugLevel,
	"info":  zerolog.InfoLevel,
	"warn":  zerolog.WarnLevel,
	"error": zerolog.ErrorLevel,
	"fatal": zerolog.FatalLevel,
}

var once sync.Once
var zeroSingleLogger *zerolog.Logger

type zeroLog struct {
	cfg    *config.Config
	logger *zerolog.Logger
}

func newZeroLog(cfg *config.Config) *zeroLog {
	logger := &zeroLog{cfg: cfg}
	logger.Init()
	return logger
}

func (l *zeroLog) getLogLevel() zerolog.Level {
	level, exists := zeroLogLevelMap[l.cfg.Logger.Level]
	if !exists {
		return zerolog.DebugLevel
	}
	return level
}

func (l *zeroLog) Init() {
	once.Do(func() {
		zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
		zerolog.SetGlobalLevel(l.getLogLevel())

		log.Printf("log file path: %s \n", getLogPath(l.cfg))
		file, err := os.OpenFile(getLogPath(l.cfg), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
		if err != nil {
			panic("could not open log file: " + err.Error())
		}
		var logger = zerolog.New(file).With().Timestamp().Str("AppName", l.cfg.App.Name).Str("LoggerName", l.cfg.Logger.Logger).Logger()
		zeroSingleLogger = &logger
	})
	l.logger = zeroSingleLogger
}

func logParamsToZeroParams(keys map[ExtraKey]any) map[string]any {
	params := make(map[string]any)
	for k, v := range keys {
		params[string(k)] = v
	}
	return params
}
func (l *zeroLog) Debug(cat Categoty, sub SubCategoty, msg string, extra map[ExtraKey]any) {
	l.logger.Debug().
		Str("Category", string(cat)).
		Str("SubCategory", string(sub)).
		Str("Extra", string(sub)).
		Fields(logParamsToZeroParams(extra)).
		Msg(msg)
}
func (l *zeroLog) Debugf(template string, args ...any) {
	l.logger.Debug().Msgf(template, args...)
}

func (l *zeroLog) Info(cat Categoty, sub SubCategoty, msg string, extra map[ExtraKey]any) {
	l.logger.Info().
		Str("Category", string(cat)).
		Str("SubCategory", string(sub)).
		Str("Extra", string(sub)).
		Fields(logParamsToZeroParams(extra)).
		Msg(msg)
}
func (l *zeroLog) Infof(template string, args ...any) {
	l.logger.Info().Msgf(template, args...)
}

func (l *zeroLog) Warn(cat Categoty, sub SubCategoty, msg string, extra map[ExtraKey]any) {
	l.logger.Warn().
		Str("Category", string(cat)).
		Str("SubCategory", string(sub)).
		Str("Extra", string(sub)).
		Fields(logParamsToZeroParams(extra)).
		Msg(msg)
}
func (l *zeroLog) Warnf(template string, args ...any) {
	l.logger.Warn().Msgf(template, args...)
}

func (l *zeroLog) Error(cat Categoty, sub SubCategoty, msg string, extra map[ExtraKey]any) {
	l.logger.Error().
		Str("Category", string(cat)).
		Str("SubCategory", string(sub)).
		Str("Extra", string(sub)).
		Fields(logParamsToZeroParams(extra)).
		Msg(msg)
}
func (l *zeroLog) Errorf(template string, args ...any) {
	l.logger.Error().Msgf(template, args...)
}

func (l *zeroLog) Fatal(cat Categoty, sub SubCategoty, msg string, extra map[ExtraKey]any) {
	l.logger.Fatal().
		Str("Category", string(cat)).
		Str("SubCategory", string(sub)).
		Str("Extra", string(sub)).
		Fields(logParamsToZeroParams(extra)).
		Msg(msg)
}
func (l *zeroLog) Fatalf(template string, args ...any) {
	l.logger.Fatal().Msgf(template, args...)
}
