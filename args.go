package logo

import "log/slog"

const (
	argComponentKey = "component"
	argErrorKey     = "error"
)

func Component(name string) slog.Attr {
	return slog.String(argComponentKey, name)
}

func Error(err error) slog.Attr {
	return slog.String(argErrorKey, err.Error())
}
