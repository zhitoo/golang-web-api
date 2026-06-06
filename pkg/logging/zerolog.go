package logging

import (
	"log"
	"os"
	"sync"

	"github.com/zhitoo/golang-web-api/config"

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

// scopedZeroLog is a zerolog-backed ScopedLogger with a fixed Category/SubCategory.
type scopedZeroLog struct {
	logger *zerolog.Logger
	cat    Category
	sub    SubCategory
}

func (l *zeroLog) With(cat Category, sub SubCategory) ScopedLogger {
	return &scopedZeroLog{logger: l.logger, cat: cat, sub: sub}
}

func (s *scopedZeroLog) logEvent(e *zerolog.Event, msg string, extra map[ExtraKey]any) {
	e.Str("Category", string(s.cat)).
		Str("SubCategory", string(s.sub)).
		Fields(logParamsToZeroParams(extra)).
		Msg(msg)
}

func (s *scopedZeroLog) Debug(msg string, extra map[ExtraKey]any) {
	s.logEvent(s.logger.Debug(), msg, extra)
}
func (s *scopedZeroLog) Debugf(template string, args ...any) {
	s.logger.Debug().Msgf(template, args...)
}
func (s *scopedZeroLog) Info(msg string, extra map[ExtraKey]any) {
	s.logEvent(s.logger.Info(), msg, extra)
}
func (s *scopedZeroLog) Infof(template string, args ...any) {
	s.logger.Info().Msgf(template, args...)
}
func (s *scopedZeroLog) Warn(msg string, extra map[ExtraKey]any) {
	s.logEvent(s.logger.Warn(), msg, extra)
}
func (s *scopedZeroLog) Warnf(template string, args ...any) {
	s.logger.Warn().Msgf(template, args...)
}
func (s *scopedZeroLog) Error(msg string, extra map[ExtraKey]any) {
	s.logEvent(s.logger.Error(), msg, extra)
}
func (s *scopedZeroLog) Errorf(template string, args ...any) {
	s.logger.Error().Msgf(template, args...)
}
func (s *scopedZeroLog) Fatal(msg string, extra map[ExtraKey]any) {
	s.logEvent(s.logger.Fatal(), msg, extra)
}
func (s *scopedZeroLog) Fatalf(template string, args ...any) {
	s.logger.Fatal().Msgf(template, args...)
}

func logParamsToZeroParams(keys map[ExtraKey]any) map[string]any {
	params := make(map[string]any)
	for k, v := range keys {
		params[string(k)] = v
	}
	return params
}
func (l *zeroLog) Debug(cat Category, sub SubCategory, msg string, extra map[ExtraKey]any) {
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

func (l *zeroLog) Info(cat Category, sub SubCategory, msg string, extra map[ExtraKey]any) {
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

func (l *zeroLog) Warn(cat Category, sub SubCategory, msg string, extra map[ExtraKey]any) {
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

func (l *zeroLog) Error(cat Category, sub SubCategory, msg string, extra map[ExtraKey]any) {
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

func (l *zeroLog) Fatal(cat Category, sub SubCategory, msg string, extra map[ExtraKey]any) {
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
