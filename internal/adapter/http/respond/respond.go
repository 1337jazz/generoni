package respond

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/1337jazz/generoni/internal/core/domain"
)

const maxBodyBytes = 1 << 20

var (
	ErrInternal      = domain.NewError("internal_error", domain.KindInternal)
	ErrMalformedBody = domain.NewError("malformed_body", domain.KindInvalidArgument)
	ErrNotFound      = domain.NewError("not_found", domain.KindNotFound)
)

type errorBody struct {
	Error errorDetail `json:"error"`
}

type errorDetail struct {
	Code string `json:"code"`
}

func Decode(w http.ResponseWriter, r *http.Request, v any) error {
	dec := json.NewDecoder(http.MaxBytesReader(w, r.Body, maxBodyBytes))
	dec.DisallowUnknownFields()
	if err := dec.Decode(v); err != nil {
		return ErrMalformedBody
	}
	return nil
}

func JSON(w http.ResponseWriter, r *http.Request, status int, v any) {
	body, err := json.Marshal(v)
	if err != nil {
		slog.ErrorContext(r.Context(), "failed to marshal response", "error", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(status)
	_, _ = w.Write(body)
}

func Error(w http.ResponseWriter, r *http.Request, err error) {
	appErr, ok := errors.AsType[*domain.Error](err)
	if !ok {
		slog.ErrorContext(r.Context(), "unexpected error", "error", err)
		appErr = ErrInternal
	}
	JSON(w, r, statusFor(appErr.Kind), errorBody{Error: errorDetail{Code: appErr.Code}})
}

func statusFor(kind domain.Kind) int {
	switch kind {
	case domain.KindInvalidArgument:
		return http.StatusBadRequest
	case domain.KindNotFound:
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}
