package server

import (
	"fmt"
	"github.com/downace/print-server/internal/logging"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"net/http"
	"net/netip"
)

func methodNotAllowed(writer http.ResponseWriter, _ *http.Request) {
	RespondError(writer, "method not allowed", http.StatusMethodNotAllowed)
}

func notFound(writer http.ResponseWriter, _ *http.Request) {
	RespondError(writer, "not found", http.StatusNotFound)
}

func panicHandlerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				RespondError(writer, fmt.Sprint(err), http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(writer, request)
	})
}

func CreateServer(addr netip.AddrPort) *http.Server {
	router := mux.NewRouter()

	router.
		Path("/printers").
		Methods("GET").
		HandlerFunc(getPrinters)

	router.
		Path("/print-pdf").
		Methods("POST").
		Headers("Content-Type", "application/json").
		HandlerFunc(printPdf)

	router.MethodNotAllowedHandler = http.HandlerFunc(methodNotAllowed)
	router.NotFoundHandler = http.HandlerFunc(notFound)

	router.Use(panicHandlerMiddleware)

	handler := handlers.CombinedLoggingHandler(logging.HttpLog.Writer(), router)

	return &http.Server{
		Addr:    addr.String(),
		Handler: handler,
	}
}
