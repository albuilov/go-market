// Package logging создает структурированный логгер приложения.
package logging

import (
	"io"
	"log/slog"
)

// New создает JSON-логгер и добавляет имя сервиса во все записи.
func New(
	output io.Writer,
	service string,
	level slog.Leveler,
) *slog.Logger {
	handler := slog.NewJSONHandler(
		output,
		&slog.HandlerOptions{Level: level},
	)

	return slog.New(handler).With(slog.String("service", service))
}
