package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	addr := flag.String("addr", ":4000", "HTTP Network Address")

	flag.Parse()
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}))

	app := App{
		logger: logger,
	}

	logger.Info("Golang server started.", "address", *addr)
	err := http.ListenAndServe(*addr, app.routes())
	logger.Error(err.Error())
	os.Exit(1)
}
