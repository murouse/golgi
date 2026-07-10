# golgi

Минималистичная, высокопроизводительная обертка над стандартным `log/slog`. Использует **Uber Zap** в качестве быстрого движка («frontend» — `slog`, «backend» — `zap`) и предоставляет удобный конвейер (Middleware) для работы с контекстом.

## Ключевые фичи

* **Производительность Zap + API slog**: Никаких компромиссов между скоростью логирования и чистотой архитектуры.
* **Контекстные атрибуты**: Автоматический проброс метаданных (`RequestID`, `TraceID`) через `context.Context`.
* **Трассировочный буфер (Process Log)**: Возможность изолированно собрать все логи конкретной транзакции в память для отправки во внешние системы отладки при возникновении ошибок.
* **Отсутствие глобального состояния**: Настройка происходит через подмену стандартного `slog.SetDefault`.

---

## Быстрый старт

### 1. Инициализация в `main.go`

```go
package main

import (
	"log/slog"
	"github.com/murouse/golgi"
)

func main() {
	// Инициализация логгера (обычно на старте приложения)
	err := golgi.Init(
		golgi.WithServiceName("rack-reservation-service"),
		golgi.WithLevel(golgi.LevelInfo),
		golgi.WithFormat(golgi.FormatJSON),
	)
	if err != nil {
		panic(err)
	}

	// Использование стандартного slog API
	slog.Info("сервис успешно запущен", golgi.Component("main"))
}

```

---

## Продвинутое использование

### Логирование контекстных данных (Middleware)

Вы можете привязать любые атрибуты к контексту запроса, и они автоматически появятся во всех последующих лог-линиях внутри этой горутины:

```go
func HandleOrder(ctx context.Context, orderID int) {
	// Прошиваем ID заказа в контекст
	ctx = golgi.WithAttrs(ctx, slog.Int("order_id", orderID))

	// В логе автоматически будет и service, и component, и order_id
	slog.InfoContext(ctx, "обработка заказа началась", golgi.Component("worker"))
}

```

### Сбор логов процесса в память (Buffer Log)

Если вам нужно собрать цепочку логов внутри одной операции (например, сложного воркера) и зафиксировать их отдельно, используйте `WithBufferLog`:

```go
func ProcessTask(ctx context.Context) {
	// Активируем буфер для этого контекста
	ctx, buf := golgi.WithBufferLog(ctx)

	slog.InfoContext(ctx, "шаг 1 выполнен")
	slog.InfoContext(ctx, "шаг 2 выполнен")

	// Если что-то пошло не так, у вас есть весь слепок логов операции
	if err := doSomething(); err != nil {
		slog.ErrorContext(ctx, "ошибка процесса", golgi.Error(err))
		
		// buf.Bytes() содержит плоский текст всех логов, прошедших через этот ctx
		saveDumpToDatabase(buf.Bytes()) 
	}
}

```

---

## Архитектура конвейера (Pipeline)

Запись лога проходит через слои матрешки изнутри наружу:

$$\text{slog.Record} \longrightarrow \underbrace{\text{ContextAttrsHandler}}_{\text{Внедряет данные из ctx}} \longrightarrow \underbrace{\text{BufferHandler}}_{\text{Дублирует в память}} \longrightarrow \underbrace{\text{zapslog.Handler}}_{\text{Мост в Zap}} \longrightarrow \text{Uber Zap Core}$$

## Конфигурация

Доступные опции для `golgi.Init()`:

| Опция | Тип | Дефолт | Описание |
| --- | --- | --- | --- |
| `WithLevel()` | `golgi.Level` | `"debug"` | `debug`, `info`, `warn`, `error` |
| `WithFormat()` | `golgi.Format` | `"json"` | `"json"` или `"console"` |
| `WithCaller()` | `bool` | `true` | Добавляет строку вызова лога в код |
| `WithServiceName()` | `string` | `""` | Добавляет поле `service` во все логи |