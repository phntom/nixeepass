package orm

import (
	"context"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm/logger"
	"time"
)

type Logger struct {
}

func (l Logger) LogMode(logger.LogLevel) logger.Interface {
	return l
}

func (l Logger) Error(ctx context.Context, msg string, opts ...interface{}) {
	log.Ctx(ctx).Error().Msg(fmt.Sprintf(msg, opts...))
}

func (l Logger) Warn(ctx context.Context, msg string, opts ...interface{}) {
	log.Ctx(ctx).Warn().Msg(fmt.Sprintf(msg, opts...))
}

func (l Logger) Info(ctx context.Context, msg string, opts ...interface{}) {
	log.Ctx(ctx).Info().Msg(fmt.Sprintf(msg, opts...))
}

func (l Logger) Trace(ctx context.Context, begin time.Time, f func() (string, int64), err error) {
	var event = log.Ctx(ctx).Trace()
	if err != nil {
		event = log.Ctx(ctx).Debug().Err(err)
	}

	sql, rows := f()
	if rows > -1 {
		event.Int64("rows", rows)
	}
	if sql != "" {
		event.Str("sql", sql)
	}

	var durKey string

	switch zerolog.DurationFieldUnit {
	case time.Nanosecond:
		durKey = "elapsed_ns"
	case time.Microsecond:
		durKey = "elapsed_us"
	case time.Millisecond:
		durKey = "elapsed_ms"
	case time.Second:
		durKey = "elapsed"
	case time.Minute:
		durKey = "elapsed_min"
	case time.Hour:
		durKey = "elapsed_hr"
	default:
		log.Error().Interface("zerolog.DurationFieldUnit", zerolog.DurationFieldUnit).Msg("gormzerolog encountered a mysterious, unknown value for DurationFieldUnit")
		durKey = "elapsed_"
	}

	event.Dur(durKey, time.Since(begin))

	event.Send()
}
