package logger

import (
	"context"

	"github.com/jackc/pgx/v5/tracelog"
	"github.com/rs/zerolog"
)

type zerologPgxLogger struct {
	logger zerolog.Logger
}

func (l zerologPgxLogger) Log(ctx context.Context, level tracelog.LogLevel, msg string, data map[string]interface{}) {
	logEvent := l.logger.With().Fields(data).Logger()
	zerologLevel := l.logger.GetLevel()
	if zerologLevel == zerolog.TraceLevel {
		level = tracelog.LogLevelTrace
	}

	switch level {
	case tracelog.LogLevelTrace:
		logEvent.Trace().Msg(msg)
	case tracelog.LogLevelDebug:
		logEvent.Debug().Msg(msg)
	case tracelog.LogLevelInfo:
		logEvent.Info().Msg(msg)
	case tracelog.LogLevelWarn:
		logEvent.Warn().Msg(msg)
	case tracelog.LogLevelError:
		logEvent.Error().Msg(msg)
	default:
		logEvent.Info().Msg(msg)
	}
}

func NewTraceLogger(zeroLog zerolog.Logger, isDebug bool) *tracelog.TraceLog {
	logLevel := tracelog.LogLevelInfo
	if isDebug {
		logLevel = tracelog.LogLevelTrace
	}

	return &tracelog.TraceLog{
		Logger:   zerologPgxLogger{logger: zeroLog},
		LogLevel: logLevel,
		Config: &tracelog.TraceLogConfig{
			TimeKey: "duration",
		},
	}
}
