package logo

import (
	"go.uber.org/zap/zapcore"
)

// levelMap транслирует строковый доменный уровень библиотеки в низкоуровневый шаг шкалы zapcore.Level.
var zapLevelMap = map[Level]zapcore.Level{
	LevelDebug: zapcore.DebugLevel,
	LevelInfo:  zapcore.InfoLevel,
	LevelWarn:  zapcore.WarnLevel,
	LevelError: zapcore.ErrorLevel,
}

// LevelToZapLevel адаптирует уровни под требования ядра Zap, предотвращая рассинхронизацию шкал.
func LevelToZapLevel(level Level) zapcore.Level {
	return zapLevelMap[level]
}
