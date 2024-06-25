package logger

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/s-larionov/process-manager"
)

type ProcessManagerLogger struct {
}

func (l *ProcessManagerLogger) Info(msg string, fields ...process.LogFields) {
	log.Info().Fields(convertFields(fields)).Msg(msg)
}

func (l *ProcessManagerLogger) Error(msg string, err error, fields ...process.LogFields) {
	Critical(err).Fields(convertFields(fields)).Msg(msg)
}

func convertFields(fields []process.LogFields) map[string]interface{} {
	if len(fields) == 0 {
		return nil
	}

	return fields[0]
}

func Critical(err error) *zerolog.Event {
	return log.Error().Err(err).Str("severity", "CRITICAL")
}
