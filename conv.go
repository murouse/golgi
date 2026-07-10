package golgi

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

// toZap адаптирует уровни под требования ядра Zap, предотвращая рассинхронизацию шкал.
func (lvl Level) toZap() zapcore.Level {
	return zapLevelMap[lvl]
}

var encodeCallerMap = map[EncodeCaller]zapcore.CallerEncoder{
	EncodeCallerShort: zapcore.ShortCallerEncoder,
	EncodeCallerFull:  zapcore.FullCallerEncoder,
	EncodeCallerSmart: SmartCallerEncoder,
}

func (ce EncodeCaller) toZap() zapcore.CallerEncoder {
	return encodeCallerMap[ce]
}
