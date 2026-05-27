package handlers

import (
	"context"
	"log/slog"

	"github.com/murouse/logo/logctx"
)

// ContextAttrsHandler автоматически обогащает каждую запись лога атрибутами,
// сохраненными ранее в контексте выполнения горутины.
type ContextAttrsHandler struct {
	slog.Handler
}

// NewContextAttrsHandler создает хендлер-обогатитель контекста.
func NewContextAttrsHandler(next slog.Handler) slog.Handler {
	return &ContextAttrsHandler{next}
}

// Handle вытаскивает атрибуты из контекста, внедряет их в slog.Record и передает дальше по цепочке.
func (h *ContextAttrsHandler) Handle(ctx context.Context, r slog.Record) error {
	r.AddAttrs(logctx.AttrsFromContext(ctx)...)
	return h.Handler.Handle(ctx, r)
}

func (h *ContextAttrsHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &ContextAttrsHandler{Handler: h.Handler.WithAttrs(attrs)}
}

func (h *ContextAttrsHandler) WithGroup(name string) slog.Handler {
	return &ContextAttrsHandler{Handler: h.Handler.WithGroup(name)}
}
