package logo

import "log/slog"

// Строковые константы для жесткой стандартизации системных ключей в JSON-файлах логов.
const (
	argServiceKey   = "service"
	argComponentKey = "component"
	argErrorKey     = "error"
)

// Component типизирует логируемый архитектурный модуль (например: "postgres", "kafka-consumer").
func Component(name string) slog.Attr {
	return slog.String(argComponentKey, name)
}

// Error приводит ошибку к стандартному ключу "error" со строковым значением (вместо сериализации структуры err).
func Error(err error) slog.Attr {
	if err == nil {
		return slog.Attr{}
	}
	return slog.String(argErrorKey, err.Error())
}

// Service типизирует логируемый сервис.
func Service(name string) slog.Attr {
	return slog.String(argServiceKey, name)
}
