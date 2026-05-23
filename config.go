package logo

import (
	"fmt"
	"io"
)

// Config объединяет все параметры кастомизации логгера.
type Config struct {
	Level        Level        // Уровень логирования (строка: debug, info, warn, error)
	Format       Format       // Формат вывода логов (строка: json, console)
	WithCaller   bool         // Флаг добавления места вызова в лог (file.go:line)
	CallerSkip   int          // CallerSkip корректирует смещение кадра стека для точного определения места вызова.
	ServiceName  *string      // Имя сервиса, сквозным образом добавляемое во все записи
	EncodeCaller EncodeCaller // EncodeCaller задает стратегию форматирования путей к файлам исходного кода.
	Writer       io.Writer    // Writer определяет целевой поток вывода логов.
	Middlewares  []Middleware
}

// Apply последовательно накатывает функциональные опции на текущую структуру конфигурации.
func (c *Config) Apply(opts ...Option) error {
	for _, opt := range opts {
		if err := opt(c); err != nil {
			return fmt.Errorf("apply option: %w", err)
		}
	}

	return nil
}

// Option инкапсулирует замыкание для безопасной настройки полей Config с валидацией "на лету".
type Option func(*Config) error

// WithLevel проверяет уровень по белому списку и прошивает его в конфиг.
func WithLevel(level Level) Option {
	return func(config *Config) error {
		_, ok := levelMap[level]
		if !ok {
			return fmt.Errorf("invalid level: %s", level)
		}
		config.Level = level
		return nil
	}
}

// WithFormat проверяет формат (json/console) и прошивает его в конфиг.
func WithFormat(format Format) Option {
	return func(config *Config) error {
		_, ok := formatMap[format]
		if !ok {
			return fmt.Errorf("invalid format: %s", format)
		}
		config.Format = format
		return nil
	}
}

// WithServiceName задает глобальный идентификатор сервиса.
func WithServiceName(serviceName string) Option {
	return func(config *Config) error {
		config.ServiceName = &serviceName
		return nil
	}
}

// WithCaller управляет отображением метаданных исходного кода в лог-линии.
func WithCaller(enabled bool) Option {
	return func(config *Config) error {
		config.WithCaller = enabled
		return nil
	}
}

// WithCallerSkip изменяет глубину обхода стека вызовов для runtime.Caller.
func WithCallerSkip(skip int) Option {
	return func(cfg *Config) error {
		cfg.CallerSkip = skip
		return nil
	}
}

// WithEncodeCaller проверяет стратегию форматирования путей и сохраняет её в конфиг.
func WithEncodeCaller(encodeCaller EncodeCaller) Option {
	return func(config *Config) error {
		_, ok := encodeCallerMap[encodeCaller]
		if !ok {
			return fmt.Errorf("invalid encode caller: %s", encodeCaller)
		}
		config.EncodeCaller = encodeCaller
		return nil
	}
}

// WithWriter перенаправляет поток вывода логов в кастомный io.Writer.
func WithWriter(writer io.Writer) Option {
	return func(config *Config) error {
		config.Writer = writer
		return nil
	}
}

func WithMiddlewares(middlewares ...Middleware) Option {
	return func(config *Config) error {
		config.Middlewares = append(config.Middlewares, middlewares...)
		return nil
	}
}

func WithResetMiddlewares(middlewares ...Middleware) Option {
	return func(config *Config) error {
		config.Middlewares = middlewares
		return nil
	}
}
