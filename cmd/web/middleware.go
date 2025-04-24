package main

import (
	"net/http"
)

func (app *application) loggingMiddleware(next http.Handler) http.Handler {
	fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			ip     = r.RemoteAddr
			proto  = r.Proto
			method = r.Method
			uri    = r.URL.RequestURI()
		)

		app.logger.Info("received request", "ip", ip, "protocol", proto, "method", method, "uri", uri)
		next.ServeHTTP(w, r)
	})
	return fn

}

func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if the user is authenticated by looking at the session
		session, err := app.sessionStore.Get(r, "session-name")
		if err != nil {
			app.logger.Error("error getting session", "error", err)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		isAuthenticated, ok := session.Values["authenticated"].(bool)
		if !ok || !isAuthenticated {
			app.logger.Info("unauthenticated user", "uri", r.URL.Path)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// If the user is authenticated, allow the request to proceed to the next handler
		next.ServeHTTP(w, r)
	})
}

func (app *application) requireRole(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session, err := app.sessionStore.Get(r, "session-name")
			if err != nil {
				app.logger.Error("error getting session for authorization", "error", err, "path", r.URL.Path)
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}

			userRole, ok := session.Values["userRole"].(string)
			if !ok || userRole != role {
				app.logger.Info("unauthorized access attempt", "path", r.URL.Path, "user_role", userRole, "required_role", role)
				app.forbidden(w,r) // Implement a forbidden handler
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

