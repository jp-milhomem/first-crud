package main

import (
	"errors"
	"log/slog"
	"net/http"
	"os"

	"github.com/jp-milhomem/first-crud/server"
)

func main() {
	slog.Info("Iniciando app - CRUD de usuários")

	if err := server.Create(); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			slog.Warn("Shutdown the server and services...")
			os.Exit(0)
		}

		slog.Error("Internal server error")
		os.Exit(1)
	}
}
