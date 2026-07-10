package golgi

import (
	"fmt"
	"os"

	"github.com/murouse/golgi/handlers"
)

// Default возвращает базовый неизменяемый пресет настроек для локальной разработки.
func Default() *Config {
	return &Config{
		Level:        LevelDebug,
		Format:       FormatJSON,
		WithCaller:   true,
		CallerSkip:   1,
		ServiceName:  nil,
		EncodeCaller: EncodeCallerSmart,
		Writer:       os.Stdout,
		Middlewares: []Middleware{
			handlers.NewContextAttrsHandler, // Порядок применения важен: ContextAttrsHandler должен отработать ДО BufferHandler, чтобы буфер увидел уже обогащенные контекстом записи.
			handlers.NewBufferHandler,
		},
	}
}

// DefaultWith создает дефолтный конфиг и сразу модифицирует его переданными опциями.
func DefaultWith(opts ...Option) (*Config, error) {
	cfg := Default()
	if err := cfg.Apply(opts...); err != nil {
		return nil, fmt.Errorf("apply options: %w", err)
	}
	return cfg, nil
}
