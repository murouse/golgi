package logo

import (
	"log/slog"
	"os"

	"github.com/murouse/logo/core"
	"github.com/murouse/logo/handlers"
	"go.uber.org/zap/exp/zapslog"
	"go.uber.org/zap/zapcore"
)

func init() {
	Init() // Дефолтная инициализация
}

// Init инициализирует вручную
func Init(opts ...Option) {
	cfg := DefaultConfig()
	for _, opt := range opts {
		opt(cfg)
	}

	// Core
	zapLogger := core.NewZapLogger(toZapLevel(cfg.Level.Level()), cfg.UseJSON)
	handler := slog.Handler(zapslog.NewHandler(zapLogger.Core(), zapslog.WithCaller(true)))

	// Middlewares
	for _, middleware := range []Middleware{
		handlers.NewContextAttrsHandler,
		handlers.NewBufferHandler,
	} {
		handler = middleware(handler)
	}

	slog.SetDefault(slog.New(handler))
}

func toZapLevel(l slog.Level) zapcore.Level {
	return zapcore.Level(l) // Простой маппинг slog.Level -> zapcore.Level
}

type Middleware func(slog.Handler) slog.Handler

func NewHandler(cfg *Config) slog.Handler {
	var handler slog.Handler
	opts := &slog.HandlerOptions{
		Level:     cfg.Level,
		AddSource: cfg.AddSource,
	}

	if cfg.UseJSON {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	return handler
}
