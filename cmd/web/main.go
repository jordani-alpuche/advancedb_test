package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"html/template"
	"log/slog"
	"os"
	"time"

	"github/jordani-alpuche/test1/internal/data"

	_ "github.com/lib/pq"
)

type application struct {
	addr *string
	brandInfo      *data.BrandDataModel
	categoryInfo	 *data.CategoryDataModel
	productInfo         *data.ProductDataModel
	logger        *slog.Logger
	templateCache map[string]*template.Template
}

func main() {
	addr := flag.String("addr", "", "HTTP network address")
	dsn := flag.String("dsn", "", "PostgreSQL DSN")

	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	db, err := openDB(*dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	logger.Info("database connection pool established")
	templateCache, err := newTemplateCache()
	logger.Info(fmt.Sprintf("template cache has %d templates", len(templateCache)))
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	defer db.Close()

	app := &application{
		addr: addr,
		brandInfo:      &data.BrandDataModel{DB: db},
		categoryInfo: &data.CategoryDataModel{DB: db},
		productInfo:         &data.ProductDataModel{DB: db},
		logger:        logger,
		templateCache: templateCache,
	}

	err = app.serve()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
