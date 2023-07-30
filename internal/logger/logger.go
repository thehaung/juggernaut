package logger

import (
	"fmt"
	"github.com/rs/zerolog"
	"os"
	"sync"
	"time"
)

const _defaultEnvironmentEnv = "APP_ENV"

var (
	loggerInstance *Logger
	syncOnce       sync.Once
)

// Interface -.
type Interface interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
}

type Logger struct {
	log zerolog.Logger
}

func GetLogger() Interface {
	syncOnce.Do(func() {
		// This is hard config for skip get caller from root execution
		skipFrameCount := 2
		zerolog.TimeFieldFormat = time.RFC3339Nano
		logger := zerolog.New(os.Stdout).
			With().
			Timestamp().
			CallerWithSkipFrameCount(zerolog.CallerSkipFrameCount + skipFrameCount).
			Logger()

		if os.Getenv(_defaultEnvironmentEnv) != "production" {
			logger = zerolog.New(os.Stdout).
				With().
				Timestamp().
				CallerWithSkipFrameCount(zerolog.CallerSkipFrameCount + 2).
				Logger().
				Output(zerolog.ConsoleWriter{Out: os.Stderr})
		}
		l := &Logger{
			log: logger,
		}
		loggerInstance = l
	})

	return loggerInstance
}

func (l *Logger) Print(level zerolog.Level, args ...interface{}) {
	l.log.WithLevel(level).Msg(fmt.Sprint(args...))
}

func (l *Logger) Printf(level zerolog.Level, format string, v ...interface{}) {
	l.log.WithLevel(zerolog.DebugLevel).Msgf(format, v...)
}

func (l *Logger) Debug(args ...interface{}) {
	l.Print(zerolog.DebugLevel, args...)
}

func (l *Logger) Info(args ...interface{}) {
	l.Print(zerolog.InfoLevel, args...)
}

func (l *Logger) Warn(args ...interface{}) {
	l.Print(zerolog.WarnLevel, args...)
}

func (l *Logger) Error(args ...interface{}) {
	l.Print(zerolog.ErrorLevel, args...)
}

func (l *Logger) Fatal(args ...interface{}) {
	l.Print(zerolog.FatalLevel, args...)
	os.Exit(1)
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	l.Printf(zerolog.DebugLevel, format, args...)
}

func (l *Logger) Infof(format string, args ...interface{}) {
	l.Printf(zerolog.InfoLevel, format, args...)
}

func (l *Logger) Warnf(format string, args ...interface{}) {
	l.Printf(zerolog.WarnLevel, format, args...)
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	l.Printf(zerolog.ErrorLevel, format, args...)
}

func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.Printf(zerolog.FatalLevel, format, args...)
	os.Exit(1)
}
