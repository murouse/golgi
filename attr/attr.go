package attr

import (
	"github.com/google/uuid"
	"log/slog"
	"runtime/debug"
)

// Строковые константы для жесткой стандартизации системных ключей в JSON-файлах логов.
const (
	attrServiceKey   = "service"
	attrComponentKey = "component"
	attrErrorKey     = "error"
	attrPanicKey     = "panic"
	attrStackKey     = "stack"
	attrJobIDKey     = "job_id"
	attrJobNameKey   = "job_name"
)

// Component типизирует логируемый архитектурный модуль (например: "postgres", "kafka-consumer").
func Component(name string) slog.Attr {
	return slog.String(attrComponentKey, name)
}

// Error приводит ошибку к стандартному ключу "error" со строковым значением (вместо сериализации структуры err).
func Error(err error) slog.Attr {
	if err == nil {
		return slog.Attr{}
	}
	return slog.String(attrErrorKey, err.Error())
}

// Service типизирует логируемый сервис.
func Service(name string) slog.Attr {
	return slog.String(attrServiceKey, name)
}

func Panic(p any) slog.Attr {
	if p == nil {
		return slog.Attr{}
	}
	return slog.Any(attrPanicKey, p)
}

func Stack() slog.Attr {
	return slog.String(attrStackKey, string(debug.Stack()))
}

func JobID(id uuid.UUID) slog.Attr {
	return slog.String(attrJobIDKey, id.String())
}

func JobName(name string) slog.Attr {
	return slog.String(attrJobNameKey, name)
}

func UUID(key string, val uuid.UUID) slog.Attr {
	return slog.String(key, val.String())
}
