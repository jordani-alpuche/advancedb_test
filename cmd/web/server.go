package main

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/gorilla/csrf"
)

func (app *application) serve() error {
	csrfKey := []byte("ZQnXOK/iAwl+wMHKrQxS1VEw+9KAZUq=")
	csrfMiddleware := csrf.Protect(
		csrfKey,
		csrf.Secure(true),
		csrf.SameSite(csrf.SameSiteDefaultMode),
		csrf.HttpOnly(false),
		csrf.Path("/"),
		
		
		csrf.TrustedOrigins([]string{
			"https://localhost:4000",
			"https://192.168.242.128:4000",
			"https://192.168.100.103:4000",
			"https://192.168.242.128",
		}),
		csrf.ErrorHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			app.logger.Error("CSRF failure",
				"method", r.Method,
				"path", r.URL.Path,
				"form_token", r.FormValue("gorilla.csrf.Token"),
				"header_token", r.Header.Get("X-CSRF-Token"),
				"cookie", r.Header.Get("Cookie"),
				"origin", r.Header.Get("Origin"),
				"referer", r.Referer(),
				)
				http.Error(w, "Forbidden - CSRF token invalid", http.StatusForbidden)
		},
		),
	),
	)

	srv := &http.Server{
		Addr:         *app.addr,
		Handler:      csrfMiddleware(app.routes()),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     slog.NewLogLogger(app.logger.Handler(), slog.LevelError),
	}
	app.logger.Info("starting server", "addr", srv.Addr, "handler", srv.Handler)
	return srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
}
