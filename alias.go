package logo

import (
	"context"
	"log/slog"

	"github.com/murouse/logo/logctx"
)

func WithAttrs(ctx context.Context, a ...slog.Attr) context.Context {
	return logctx.WithAttrs(ctx, a...)
}
