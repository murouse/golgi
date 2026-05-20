package handlers

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/murouse/logo/logctx"
)

type Byter interface {
	Bytes() []byte
}

type safeBuffer struct {
	sync.Mutex
	bytes.Buffer
}

func WithBufferLog(ctx context.Context) (context.Context, Byter) {
	buf := &safeBuffer{}
	ctx = context.WithValue(ctx, bufferLogKey, buf)
	return ctx, buf
}

var bufferLogKey struct{}

type BufferHandler struct {
	slog.Handler
}

func NewBufferHandler(next slog.Handler) slog.Handler {
	return &BufferHandler{next}
}

func (h *BufferHandler) Handle(ctx context.Context, rec slog.Record) error {
	// пишем в основной лог
	err := h.Handler.Handle(ctx, rec)

	// дублируем в process-log, если он есть
	if buf, ok := ctx.Value(bufferLogKey).(*safeBuffer); ok {
		buf.Lock()
		defer buf.Unlock()

		ts := rec.Time.Format(time.RFC3339)

		fmt.Fprintf(buf, "[%s] %s - %s", ts, rec.Level, rec.Message)

		// атрибуты из контекста (помещаются в атрибуты лога только в next хендлере)
		for _, a := range logctx.AttrsFromContext(ctx) {
			fmt.Fprintf(buf, " [%s=%s]", a.Key, a.Value)
		}

		// атрибуты конкретно этого лога
		rec.Attrs(func(a slog.Attr) bool {
			fmt.Fprintf(buf, " [%s=%s]", a.Key, a.Value)
			return true
		})

		buf.WriteByte('\n') // nolint: errcheck
	}

	return err
}
