package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

type MalformedRequest struct {
	Status  int
	message string
}

func (mr MalformedRequest) Error() string {
	return mr.message
}

func DecodeJSONBody(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	if ct := r.Header.Get("Content-Type"); ct != "" {
		mediaType := strings.ToLower(strings.TrimSpace(strings.Split(ct, ";")[0]))
		if mediaType != "application/json" {
			return &MalformedRequest{
				Status:  http.StatusUnsupportedMediaType,
				message: "Content Type header is not application/json",
			}
		}
	}

	// Enforce read only 1MB from request body
	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(&dst); err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxError):
			return &MalformedRequest{
				Status: http.StatusBadRequest,
				message: fmt.Sprintf(
					"Request body contains badly-formed JSON (at position %d)",
					syntaxError.Offset,
				),
			}

		case errors.Is(err, io.ErrUnexpectedEOF):
			return &MalformedRequest{Status: http.StatusBadRequest, message: "Request body contains badly-formed JSON"}

		case errors.As(err, &unmarshalTypeError):
			msg := fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
			return &MalformedRequest{
				Status:  http.StatusBadRequest,
				message: msg,
			}

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			msg := fmt.Sprintf("Request body contains unknown field %s", fieldName)
			return &MalformedRequest{
				Status:  http.StatusBadRequest,
				message: msg,
			}

		case errors.Is(err, io.EOF):
			return &MalformedRequest{
				Status:  http.StatusBadRequest,
				message: "Request body must not be empty",
			}

		case err.Error() == "http: request body too large":
			return &MalformedRequest{
				Status:  http.StatusRequestEntityTooLarge,
				message: "Request body must not be larger than 1MB",
			}

		default:
			log.Print(err.Error())
			return &MalformedRequest{
				Status:  http.StatusInternalServerError,
				message: http.StatusText(http.StatusInternalServerError),
			}
		}
	}

	return nil
}
