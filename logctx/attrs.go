package logctx

import (
	"context"
	"log/slog"
)

type attrsCtxKey struct{}

func WithAttrs(ctx context.Context, a ...slog.Attr) context.Context {
	attrs := append(AttrsFromContext(ctx), a...)
	return context.WithValue(ctx, attrsCtxKey{}, attrs)
}

func AttrsFromContext(ctx context.Context) []slog.Attr {
	attrs, _ := ctx.Value(attrsCtxKey{}).([]slog.Attr)
	return attrs
}
