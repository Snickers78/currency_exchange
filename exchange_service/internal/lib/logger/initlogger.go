package logg

import (
	"log/slog"
	"os"
)

func InitLogger(env string) *slog.Logger {
	var log *slog.Logger

	file, err := os.OpenFile("logs.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Error("Не удалось открыть файл для логирования", err)
	}

	switch env {
	case "local":
		log = slog.New(
			slog.NewTextHandler(file, &slog.HandlerOptions{Level: slog.LevelDebug, AddSource: true}),
		)
	case "dev":
		log = slog.New(
			slog.NewJSONHandler(file, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case "prod":
		log = slog.New(
			slog.NewJSONHandler(file, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
