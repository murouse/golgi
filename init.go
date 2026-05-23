package logo

import (
	"fmt"
	"log/slog"

	"github.com/murouse/logo/attr"
	"go.uber.org/zap/exp/zapslog"
)

func init() {
	// Автоматическая базовая инициализация при импорте пакета.
	// Гарантирует, что slog не запаникует и будет писать в понятном формате даже до явного вызова Init().
	_ = Init()
}

// Init выполняет ручную оркестрацию: собирает низкоуровневый Zap, оборачивает его в zapslog.Handler,
// строит упорядоченный конвейер обработки и делает его дефолтным для всего Go-приложения.
func Init(opts ...Option) error {
	cfg, err := DefaultWith(opts...)
	if err != nil {
		return fmt.Errorf("default with: %w", err)
	}

	zapLogger := NewZapLogger(cfg.Level, cfg.Format, cfg.EncodeCaller, cfg.Writer) // Создаем производительный фундамент (Zap)

	baseHandler := zapslog.NewHandler(
		zapLogger.Core(),
		zapslog.WithCaller(cfg.WithCaller),
		zapslog.WithCallerSkip(cfg.CallerSkip),
	) // Создаем адаптер-мост из zap в стандартный интерфейс slog.Handler

	handler := slog.Handler(baseHandler) // Приведение к интерфейсу slog.Handler необходимо, чтобыMiddleware-обертки могли прозрачно мутировать типы

	// Если задано имя сервиса, пришиваем его к базовому хендлеру
	if cfg.ServiceName != nil {
		handler = handler.WithAttrs([]slog.Attr{attr.Service(*cfg.ServiceName)})
	}

	// Собираем декораторы (Middleware) по принципу Матрешки (внутри -> наружу).
	for _, middleware := range cfg.Middlewares {
		handler = middleware(handler)
	}

	// Инжектим собранный пайплайн в стандартную библиотеку Go в качестве глобального логгера
	slog.SetDefault(slog.New(handler))
	return nil
}

// Middleware описывает контракт обертки над хендлером slog, позволяя строить цепочки (Pipeline pattern).
type Middleware func(slog.Handler) slog.Handler
