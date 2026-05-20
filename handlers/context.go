package handlers

import (
	"context"
	"log/slog"

	"github.com/murouse/logo/logctx"
)

type ContextAttrsHandler struct {
	slog.Handler
}

func NewContextAttrsHandler(next slog.Handler) slog.Handler {
	return &ContextAttrsHandler{next}
}

func (h *ContextAttrsHandler) Handle(ctx context.Context, r slog.Record) error {
	r.AddAttrs(logctx.AttrsFromContext(ctx)...)
	return h.Handler.Handle(ctx, r)
}
