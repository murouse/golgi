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

// Byter предоставляет потокобезопасный доступ к накопленным байтам лога в рамках одной бизнес-операции.
type Byter interface {
	Bytes() []byte
}

// safeBuffer защищает bytes.Buffer мьютексом для бесконфликтной параллельной записи из разных горутин.
type safeBuffer struct {
	sync.Mutex
	bytes.Buffer
}

// WithBufferLog активирует "запись в память" для текущего контекста.
// Все последующие логи процесса будут дублироваться в возвращаемый Byter.
func WithBufferLog(ctx context.Context) (context.Context, Byter) {
	buf := &safeBuffer{}
	ctx = context.WithValue(ctx, bufferCtxKey, buf)
	return ctx, buf
}

// bufferCtxKey предотвращает коллизии ключей контекста между пакетами
var bufferCtxKey struct{}

// BufferHandler — декоратор, дублирующий структурированные логи приложения
// в текстовый буфер контекста (полезно для отладки конкретных транзакций/запросов).
type BufferHandler struct {
	slog.Handler
}

// NewBufferHandler инициализирует слой перехвата логов в буфер.
func NewBufferHandler(next slog.Handler) slog.Handler {
	return &BufferHandler{next}
}

// Handle отправляет запись в основной логгер и, при наличии буфера в контексте, форматирует лог в плоский текст.
func (h *BufferHandler) Handle(ctx context.Context, rec slog.Record) error {
	err := h.Handler.Handle(ctx, rec) // Первым делом пишем лог в основной поток (Stdout)

	// Проверяем, включено ли трассировочное логирование в буфер для этого контекста
	if buf, ok := ctx.Value(bufferCtxKey).(*safeBuffer); ok {
		buf.Lock()
		defer buf.Unlock()

		ts := rec.Time.Format(time.RFC3339)

		fmt.Fprintf(buf, "[%s] %s - %s", ts, rec.Level, rec.Message)

		// Извлекаем и форматируем накопленный контекст (RequestID, UserID и т.д.)
		for _, a := range logctx.AttrsFromContext(ctx) {
			fmt.Fprintf(buf, " [%s=%s]", a.Key, a.Value)
		}

		// Извлекаем локальные атрибуты текущей строки логирования
		rec.Attrs(func(a slog.Attr) bool {
			fmt.Fprintf(buf, " [%s=%s]", a.Key, a.Value)
			return true
		})

		buf.WriteByte('\n')
	}

	return err
}

func (h *BufferHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &BufferHandler{Handler: h.Handler.WithAttrs(attrs)}
}

func (h *BufferHandler) WithGroup(name string) slog.Handler {
	return &BufferHandler{Handler: h.Handler.WithGroup(name)}
}
