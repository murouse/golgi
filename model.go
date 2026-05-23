package logo

import (
	"log/slog"
)

// Level абстрагирует строковое представление уровней от внутренней шкалы slog/zap.
type Level string

const (
	LevelDebug Level = "debug"
	LevelInfo  Level = "info"
	LevelWarn  Level = "warn"
	LevelError Level = "error"
)

// levelMap связывает кастомные строковые уровни со стандартными типами slog.Leveler.
var levelMap = map[Level]slog.Leveler{
	LevelDebug: slog.LevelDebug,
	LevelInfo:  slog.LevelInfo,
	LevelWarn:  slog.LevelWarn,
	LevelError: slog.LevelError,
}

// Format строго типизирует поддерживаемые кодировщики логов.
type Format string

const (
	FormatJSON    Format = "json"
	FormatConsole Format = "console"
)

// formatMap используется для быстрой O(1) проверки валидности формата при инициализации.
var formatMap = map[Format]struct{}{
	FormatJSON:    {},
	FormatConsole: {},
}

type EncodeCaller string

const (
	EncodeCallerShort EncodeCaller = "short"
	EncodeCallerFull  EncodeCaller = "full"
	EncodeCallerSmart EncodeCaller = "smart"
)
