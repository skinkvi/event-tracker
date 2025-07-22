package logger

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/fatih/color"
)

type CustomHandler struct {
	handler slog.Handler
}

func NewCustomHandler(opts *slog.HandlerOptions) *CustomHandler {
	return &CustomHandler{
		handler: slog.NewTextHandler(nil, opts),
	}
}

func (h *CustomHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

func (h *CustomHandler) Handle(ctx context.Context, record slog.Record) error {
	timeStr := record.Time.Format("02-01-2006 15:04")
	levelStr := color.GreenString(record.Level.String())
	msg := record.Message

	output := fmt.Sprintf("time=%s level=%s msg=%s\n", timeStr, levelStr, msg)

	record.Attrs(func(a slog.Attr) bool {
		output += fmt.Sprintf(" %s=%v", a.Key, a.Value)
		return true
	})

	if output != fmt.Sprintf("time=%s level=%s msg=%s\n", timeStr, levelStr, msg) {
		output += "\n"
	}

	if record.Level >= slog.LevelError {
		_, err := color.New(color.FgRed).Println(output)
		return err
	}

	_, err := color.New(color.FgBlue).Println(output)
	return err
}

func (h *CustomHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &CustomHandler{handler: h.handler.WithAttrs(attrs)}
}

func (h *CustomHandler) WithGroup(name string) slog.Handler {
	return &CustomHandler{handler: h.handler.WithGroup(name)}
}
