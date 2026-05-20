package logctx

import (
	"context"
	"log/slog"
)

// attrsCtxKey предотвращает коллизии ключей контекста между пакетами
type attrsCtxKey struct{}

// WithAttrs сохраняет структурированные атрибуты в контекст, склеивая их с уже существующими.
func WithAttrs(ctx context.Context, a ...slog.Attr) context.Context {
	attrs := append(AttrsFromContext(ctx), a...)
	return context.WithValue(ctx, attrsCtxKey{}, attrs)
}

// AttrsFromContext безопасно возвращает слайс атрибутов лога, изолированных в контексте.
func AttrsFromContext(ctx context.Context) []slog.Attr {
	attrs, _ := ctx.Value(attrsCtxKey{}).([]slog.Attr)
	return attrs
}
