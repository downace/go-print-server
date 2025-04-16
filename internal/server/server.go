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

func responseHeadersMiddleware(headers map[string]string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			for name, value := range headers {
				writer.Header().Set(name, value)
			}
			next.ServeHTTP(writer, request)
		})
	}
}

func CreateServer(addr netip.AddrPort, responseHeaders map[string]string) *http.Server {
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

	router.
		Path("/print-pdf-url").
		Methods("POST").
		HandlerFunc(printPdfFromUrl)

	router.MethodNotAllowedHandler = http.HandlerFunc(methodNotAllowed)
	router.NotFoundHandler = http.HandlerFunc(notFound)

	router.Use(panicHandlerMiddleware)
	router.Use(responseHeadersMiddleware(responseHeaders))

	handler := handlers.CombinedLoggingHandler(logging.HttpLog.Writer(), router)

	return &http.Server{
		Addr:    addr.String(),
		Handler: handler,
	}
}
