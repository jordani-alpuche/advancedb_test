package main

import (
	"context"
	"database/sql"
	"flag"
	"html/template"
	"log/slog"
	"os"
	"time"

	"github/jordani-alpuche/test2/internal/data"

	"github.com/gorilla/sessions"

	_ "github.com/lib/pq"
)

type application struct {
	addr *string
	brandInfo      *data.BrandDataModel
	categoryInfo	 *data.CategoryDataModel
	productInfo         *data.ProductDataModel
	userInfo 		*data.UsersDataModel
	sessionStore  *sessions.CookieStore
	logger        *slog.Logger
	templateCache map[string]*template.Template
}

func main() {
	addr := flag.String("addr", "", "HTTP network address")
	dsn := flag.String("dsn", "", "PostgreSQL DSN")
	secret := flag.String("secret", "3PJXnyuVxcfLdaZ92rae7S8", "Secret key for session cookies")

	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	db, err := openDB(*dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	logger.Info("database connection pool established")
	templateCache, err := newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	defer db.Close()

	// Initialize the session store with the provided secret
	sessionStore := sessions.NewCookieStore([]byte(*secret))
	sessionStore.Options = &sessions.Options{
		Secure: true,
		MaxAge: int(12 * time.Hour),
	}


	app := &application{	
		addr: addr,
		brandInfo:      &data.BrandDataModel{DB: db},
		categoryInfo: &data.CategoryDataModel{DB: db},
		productInfo:         &data.ProductDataModel{DB: db},
		userInfo: 		&data.UsersDataModel{DB: db},
		sessionStore:  sessionStore,
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
