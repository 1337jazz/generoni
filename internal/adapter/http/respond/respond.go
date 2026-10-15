package respond

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"mime"
	"net/http"

	"github.com/1337jazz/generoni/internal/core/domain"
)

const (
	maxBodyBytes = 1 << 20

	kindUnsupportedMediaType domain.Kind = "unsupported_media_type"
)

var (
	ErrInternal             = domain.NewError("internal_error", domain.KindInternal)
	ErrMalformedBody        = domain.NewError("malformed_body", domain.KindInvalidArgument)
	ErrBodyTooLarge         = domain.NewError("body_too_large", domain.KindTooLarge)
	ErrUnsupportedMediaType = domain.NewError("unsupported_media_type", kindUnsupportedMediaType)
	ErrTimeout              = domain.NewError("timeout", domain.KindUnavailable)
)

type errorBody struct {
	Error errorDetail `json:"error"`
}

type errorDetail struct {
	Code string `json:"code"`
}

func Decode(w http.ResponseWriter, r *http.Request, v any) error {
	if mt, _, _ := mime.ParseMediaType(r.Header.Get("Content-Type")); mt != "application/json" {
		return ErrUnsupportedMediaType
	}
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, maxBodyBytes)).Decode(v); err != nil {
		if _, ok := errors.AsType[*http.MaxBytesError](err); ok {
			return ErrBodyTooLarge
		}
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
	switch {
	case ok:
	case errors.Is(r.Context().Err(), context.DeadlineExceeded):
		slog.WarnContext(r.Context(), "request timed out", "error", err)
		appErr = ErrTimeout
	default:
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
	case domain.KindTooLarge:
		return http.StatusRequestEntityTooLarge
	case kindUnsupportedMediaType:
		return http.StatusUnsupportedMediaType
	case domain.KindUnavailable:
		return http.StatusServiceUnavailable
	default:
		return http.StatusInternalServerError
	}
}
