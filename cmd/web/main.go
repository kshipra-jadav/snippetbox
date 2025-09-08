package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/kshipra-jadav/snippetbox/internal/models"
)

func main() {
	addr := flag.String("addr", ":4000", "HTTP Network Address")

	flag.Parse()
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}))

	db, err := openDB("root:root@/snippetbox?parseTime=true")
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()

	app := App{
		logger:   logger,
		snippets: &models.SnippetsModel{DB: db},
	}

	logger.Info("Golang server started.", "address", *addr)
	err = http.ListenAndServe(*addr, app.routes())
	logger.Error(err.Error())
	os.Exit(1)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}
	fmt.Println("DB Connected")
	return db, nil
}
