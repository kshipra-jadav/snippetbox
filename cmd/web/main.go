package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
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

	templateCache, err := cacheNewTemplate()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	decoder := form.NewDecoder()

	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour

	app := App{
		logger:         logger,
		snippets:       &models.SnippetsModel{DB: db},
		templateCache:  templateCache,
		formDecoder:    decoder,
		sessionManager: sessionManager,
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
