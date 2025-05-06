package server

import (
	"encoding/json"
	"errors"
	"github.com/downace/print-server/internal/printing"
	"github.com/go-playground/form/v4"
	"github.com/go-playground/validator/v10"
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

func handleError(err error, w http.ResponseWriter) {
	if errors.Is(err, printing.ErrNotSupported) {
		RespondError(w, err.Error(), http.StatusNotImplemented)
	} else if errors.Is(err, printing.ErrRequestError) {
		RespondError(w, err.Error(), http.StatusUnprocessableEntity)
	} else {
		RespondError(w, err.Error(), http.StatusInternalServerError)
	}
}

func validateRequest[T any](r *http.Request) (*T, error) {
	dec := form.NewDecoder()
	var result T
	err := dec.Decode(&result, r.URL.Query())

	if err != nil {
		return nil, err
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	err = validate.Struct(result)

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func handleValidateRequestError(w http.ResponseWriter, err error) {
	var valErr validator.ValidationErrors
	if errors.As(err, &valErr) {
		RespondError(w, err.Error(), http.StatusUnprocessableEntity)
	} else {
		RespondError(w, err.Error(), http.StatusInternalServerError)
	}
}

func getPrinters(w http.ResponseWriter, _ *http.Request) {
	printers, err := printing.ListPrinters()

	if err != nil {
		handleError(err, w)
		return
	}

	RespondOk(w, map[string][]printing.Printer{"printers": printers})
}

type PrintPdfQuery struct {
	Printer string `form:"printer" validate:"required"`
}

func printPdf(w http.ResponseWriter, r *http.Request) {
	q, err := validateRequest[PrintPdfQuery](r)

	if err != nil {
		handleValidateRequestError(w, err)
		return
	}

	err = printing.PrintPDF(q.Printer, r.Body)

	if err != nil {
		handleError(err, w)
		return
	}

	RespondOk(w, nil)
}

type PrintPdfFromUrlQuery struct {
	Printer string `form:"printer" validate:"required"`
	Url     string `form:"url" validate:"required,url"`
}

func printPdfFromUrl(w http.ResponseWriter, r *http.Request) {
	q, err := validateRequest[PrintPdfFromUrlQuery](r)

	if err != nil {
		handleValidateRequestError(w, err)
		return
	}

	err = printing.PrintPDFFromUrl(q.Printer, q.Url)

	if err != nil {
		handleError(err, w)
		return
	}

	RespondOk(w, nil)
}

type PrintFromUrlQuery struct {
	Printer string `form:"printer" validate:"required"`
	Url     string `form:"url" validate:"required,url"`
}

func printFromUrl(w http.ResponseWriter, r *http.Request) {
	q, err := validateRequest[PrintFromUrlQuery](r)

	if err != nil {
		handleValidateRequestError(w, err)
		return
	}

	err = printing.PrintFromUrl(q.Printer, q.Url)

	if err != nil {
		handleError(err, w)
		return
	}

	RespondOk(w, nil)
}
