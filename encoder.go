package golgi

import (
	"fmt"
	"os"
	"path/filepath"

	"go.uber.org/zap/zapcore"
)

func SmartCallerEncoder(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	wd, _ := os.Getwd()

	if !caller.Defined {
		enc.AppendString("undefined")
		return
	}

	filePath := caller.File // По умолчанию берем полный путь к файлу

	if wd != "" { // Если удалось вычислить относительный путь от папки запуска - берем его
		if rel, err := filepath.Rel(wd, caller.File); err == nil {
			filePath = rel
		}
	}

	enc.AppendString(fmt.Sprintf("%s:%d", filePath, caller.Line)) // Форматируем как file.go:line
}
