package controller

import (
	"net/http"
	"time"
)

// LoggingMiddleware logs API queries
func (controller *PickUpPointController) LoggingMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		err := controller.Sender.sendAsyncMessage(LoggingMessage{
			Method: req.Method,
			URI:    req.RequestURI,
			Time:   time.Now(),
		})
		if err != nil {
			http.Error(w, "cannot send the logging message", http.StatusInternalServerError)
			return
		}

		handler.ServeHTTP(w, req)
	})
}

// AuthMiddleware authenticate user
func AuthMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		username, password, ok := req.BasicAuth()
		if !ok || username != configUsername || password != configPassword {
			http.Error(w, "authentication failed", http.StatusUnauthorized)
			return
		}
		handler.ServeHTTP(w, req)
	})
}
