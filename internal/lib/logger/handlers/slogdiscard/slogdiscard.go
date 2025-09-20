// Своя реализация хэндлера логера, в которой все методы slog.Handler интерфейса просто игнорируют все передаваемые значения и возвращают ничего не значащий результат
package slogdiscard

import (
	"context"
	"log/slog"
)

func NewDiscardLogger() *slog.Logger {
	return slog.New(NewDiscardHandler())
}

type DiscardHandler struct{}

func NewDiscardHandler() *DiscardHandler {
	return &DiscardHandler{}
}

func (d *DiscardHandler) Enabled(_ context.Context, _ slog.Level) bool {
	return false
}

func (d *DiscardHandler) Handle(_ context.Context, _ slog.Record) error {
	return nil
}

func (d *DiscardHandler) WithAttrs(_ []slog.Attr) slog.Handler {
	return d
}

func (d *DiscardHandler) WithGroup(_ string) slog.Handler {
	return d
}
