package golgi

import (
	"io"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewZapLogger создает сырой экземпляр *zap.Logger, оптимизированный под JSON или Console вывод в Stdout.
// Этот инстанс не используется напрямую в бизнес-логике, а передается как Backend для slog.
func NewZapLogger(level Level, format Format, encodeCaller EncodeCaller, writer io.Writer) *zap.Logger {
	cfg := zap.NewProductionEncoderConfig()
	cfg.EncodeTime = zapcore.ISO8601TimeEncoder // Стандартизируем таймштампы (ISO8601)

	var encoder zapcore.Encoder
	switch format {
	case FormatJSON:
		encoder = zapcore.NewJSONEncoder(cfg)
	case FormatConsole:
		cfg.EncodeLevel = zapcore.CapitalColorLevelEncoder // Расцвечивает уровни (INFO, ERROR) в терминале

		cfg.EncodeCaller = encodeCaller.toZap()

		encoder = zapcore.NewConsoleEncoder(cfg)
	}

	// Собираем ядро с прямой записью в Stdout
	core := zapcore.NewCore(encoder, zapcore.AddSync(writer), level.toZap())

	return zap.New(core)
}
