package server

import (
	"encoding/json"
	"errors"
	"github.com/downace/print-server/internal/printing"
	"net/http"
)

func RespondOk(w http.ResponseWriter, data interface{}) {
	respondJson(w, data, http.StatusOK)
}

func RespondError(w http.ResponseWriter, message string, status int) {
	respondJson(w, map[string]string{"message": message}, status)
}

func respondJson(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	body, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	_, err = w.Write(body)

	if err != nil {
		panic(err)
	}
}

func getPrinters(w http.ResponseWriter, _ *http.Request) {
	printers, err := printing.ListPrinters()

	if errors.Is(err, printing.ErrNotSupported) {
		RespondError(w, err.Error(), http.StatusNotImplemented)
		return
	}

	if err != nil {
		RespondError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	RespondOk(w, map[string][]printing.Printer{"printers": printers})
}

func printPdf(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	if !q.Has("printer") {
		RespondError(w, "printer param missing", http.StatusUnprocessableEntity)
		return
	}

	err := printing.PrintPDF(q.Get("printer"), r.Body)

	if errors.Is(err, printing.ErrNotSupported) {
		RespondError(w, err.Error(), http.StatusNotImplemented)
		return
	}

	if err != nil {
		RespondError(w, err.Error(), http.StatusInternalServerError)
	}
	RespondOk(w, nil)
}
