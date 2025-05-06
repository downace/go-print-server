package server

import (
	"encoding/base64"
	"fmt"
	"github.com/downace/print-server/internal/appconfig"
	"github.com/downace/print-server/internal/logging"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"net/http"
	"net/netip"
	"strings"
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

func checkBasicAuth(request *http.Request, username string, password string) (authorized bool) {
	authHeader := request.Header.Get("Authorization")
	if authHeader == "" {
		return false
	}
	authDataEncoded := strings.TrimSpace(strings.TrimPrefix(authHeader, "Basic "))
	authData, err := base64.StdEncoding.DecodeString(authDataEncoded)
	if err != nil {
		return false
	}
	credentials := strings.SplitN(string(authData), ":", 2)

	if len(credentials) != 2 || credentials[0] != username || credentials[1] != password {
		return false
	}

	return true
}

func basicAuthMiddleware(username string, password string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			if !checkBasicAuth(request, username, password) {
				writer.Header().Set("WWW-Authenticate", "Basic")
				RespondError(writer, "Unauthorized", http.StatusUnauthorized)
			} else {
				next.ServeHTTP(writer, request)
			}
		})
	}
}

func CreateServer(config appconfig.AppConfig) *http.Server {
	return createServer(
		netip.AddrPortFrom(netip.MustParseAddr(config.Host), config.Port),
		config.ResponseHeaders,
		config.Auth.Enabled,
		config.Auth.Username,
		config.Auth.Password,
	)
}

func createServer(
	addr netip.AddrPort,
	responseHeaders map[string]string,
	authEnabled bool,
	authUsername string,
	authPassword string,
) *http.Server {
	router := mux.NewRouter()

	router.
		Path("/printers").
		Methods("GET").
		HandlerFunc(getPrinters)

	router.
		Path("/print-pdf").
		Methods("POST").
		Headers("Content-Type", "application/pdf").
		HandlerFunc(printPdf)

	router.
		Path("/print-pdf-url").
		Methods("POST").
		HandlerFunc(printPdfFromUrl)

	router.
		Path("/print-url").
		Methods("POST").
		HandlerFunc(printFromUrl)

	router.MethodNotAllowedHandler = http.HandlerFunc(methodNotAllowed)
	router.NotFoundHandler = http.HandlerFunc(notFound)

	router.Use(panicHandlerMiddleware)
	router.Use(responseHeadersMiddleware(responseHeaders))
	if authEnabled {
		router.Use(basicAuthMiddleware(authUsername, authPassword))
	}

	handler := handlers.CombinedLoggingHandler(logging.HttpLog.Writer(), router)

	return &http.Server{
		Addr:    addr.String(),
		Handler: handler,
	}
}

func RunServer(server *http.Server, config appconfig.AppConfig) error {
	if config.TLS.Enabled {
		return server.ListenAndServeTLS(config.TLS.CertFile, config.TLS.KeyFile)
	} else {
		return server.ListenAndServe()
	}
}
