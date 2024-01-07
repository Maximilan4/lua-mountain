package slogan

import (
	"context"
	"log/slog"
	"time"
)

type (
	ContextKeyGetter interface {
		GetKey() string
	}

	commonHandler[Handler slog.Handler] struct {
		h Handler
		ContextKeys []ContextKeyGetter
	}
)


func (ch *commonHandler[Handler]) Enabled(ctx context.Context, level slog.Level) bool {
	return ch.h.Enabled(ctx, level)
}

func (ch *commonHandler[Handler]) Handle(ctx context.Context, record slog.Record) error {
	attrs := make([]slog.Attr, 0, len(ch.ContextKeys))
	for _, ck := range ch.ContextKeys {
		value := ctx.Value(ck)
		if value == nil {
			continue
		}

		//TODO: think about pointer values
		switch value.(type) {
		case string:
			attrs = append(attrs, slog.String(ck.GetKey(), value.(string)))
		case int:
			attrs = append(attrs, slog.Int(ck.GetKey(), value.(int)))
		case int64:
			attrs = append(attrs, slog.Int64(ck.GetKey(), value.(int64)))
		case uint64:
			attrs = append(attrs, slog.Uint64(ck.GetKey(), value.(uint64)))
		case float64:
			attrs = append(attrs, slog.Float64(ck.GetKey(), value.(float64)))
		case bool:
			attrs = append(attrs, slog.Bool(ck.GetKey(), value.(bool)))
		case time.Duration:
			attrs = append(attrs, slog.Duration(ck.GetKey(), value.(time.Duration)))
		case time.Time:
			attrs = append(attrs, slog.Time(ck.GetKey(), value.(time.Time)))
		case slog.LogValuer:
			attrs = append(attrs, slog.Attr{Key: ck.GetKey(), Value: value.(slog.LogValuer).LogValue()})
		default:
			attrs = append(attrs, slog.Any(ck.GetKey(), value))
		}
	}

	record.AddAttrs(attrs...)
	return ch.h.Handle(ctx, record)
}

func (ch *commonHandler[Handler]) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &commonHandler[Handler]{
		h:           ch.h.WithAttrs(attrs).(Handler),
		ContextKeys: ch.ContextKeys,
	}
}

func (ch *commonHandler[Handler]) WithGroup(name string) slog.Handler {
	return &commonHandler[Handler]{
		h:           ch.h.WithGroup(name).(Handler),
		ContextKeys: ch.ContextKeys,
	}
}


