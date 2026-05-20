package logo

import (
	"context"
	"log/slog"

	"github.com/murouse/logo/logctx"
)

// WithAttrs — публичный прокси-метод для удобного обогащения контекста логами на уровне приложения.
func WithAttrs(ctx context.Context, a ...slog.Attr) context.Context {
	return logctx.WithAttrs(ctx, a...)
}
