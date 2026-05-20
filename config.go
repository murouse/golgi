package logo

import "log/slog"

type Config struct {
	Level     slog.Leveler
	UseJSON   bool
	AddSource bool
}

func DefaultConfig() *Config {
	return &Config{
		Level:     slog.LevelDebug,
		UseJSON:   true,
		AddSource: true,
	}
}

type Option func(*Config)

func WithJSONFormat(useJSON bool) Option {
	return func(config *Config) {
		config.UseJSON = useJSON
	}
}
