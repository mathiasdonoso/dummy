package main

import (
	"log/slog"
	"os"
	"strconv"

	"github.com/mathiasdonoso/dummy/internal/cli"
)

func initLogger() {
	d, _ := strconv.Atoi(os.Getenv("DEBUG"))
	if d != 1 {
		return
	}

	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	slog.SetDefault(slog.New(handler))

	slog.Debug("debug mode enabled")
}

func main() {
	initLogger()

	if err := cli.Dispatch(os.Args[1:]); err != nil {
		slog.Debug(err.Error())
		panic(err)
	}
}
