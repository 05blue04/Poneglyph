package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/05blue04/Poneglyph/internal/validator"
	"github.com/julienschmidt/httprouter"
)

type envelope map[string]any

func (app *application) writeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
	js, err := json.MarshalIndent(data, "", "\t") //change to regular marshal in production as adding indents has an impact on performance
	if err != nil {
		return err
	}

	// Append a newline to make it easier to view in terminal
	js = append(js, '\n')

	for key, values := range headers {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	r.Body = http.MaxBytesReader(w, r.Body, 1_048_576)

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields() //prevent mfs from adding non existent parameters in their requests
	err := decoder.Decode(dst)
	if err != nil {

		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError
		var maxBytesError *http.MaxBytesError

		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)

		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")

		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)

		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")

		// A json.InvalidUnmarshalError error will be returned if we pass something
		// that is not a non-nil pointer as the target destination to Decode(). In this case we panic
		case errors.As(err, &invalidUnmarshalError):
			panic(err)

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("body contains unknown key %s", fieldName)

		case errors.As(err, &maxBytesError):
			return fmt.Errorf("body must not be larger than %d bytes", maxBytesError.Limit)
		default:
			return err
		}
	}

	err = decoder.Decode(&struct{}{}) //call decoder again-> decoder is meant to read a stream of json so if more is provided it can keep reading since we only expect one object we can raise an error if more is provided
	if !errors.Is(err, io.EOF) {
		return errors.New("body must only contain a single JSON value")
	}

	return nil
}

func parseConfig() (config, error) {
	var cfg config
	var err error

	portStr := getEnvWithDefault("PORT", "4000")
	cfg.port, err = strconv.Atoi(portStr)
	if err != nil {
		return cfg, err
	}

	cfg.env = getEnvWithDefault("ENV", "development")

	cfg.db.dsn = os.Getenv("DSN")
	if cfg.db.dsn == "" {
		return cfg, fmt.Errorf("DSN environment variable is required")
	}

	maxOpenConnsStr := getEnvWithDefault("MAX_OPEN_CONNS", "25")
	cfg.db.maxOpenConns, err = strconv.Atoi(maxOpenConnsStr)
	if err != nil {
		return cfg, err
	}

	maxIdleConnsStr := getEnvWithDefault("MAX_IDLE_CONNS", "25")
	cfg.db.maxIdleConns, err = strconv.Atoi(maxIdleConnsStr)
	if err != nil {
		return cfg, err
	}

	maxIdleTimeStr := getEnvWithDefault("MAX_IDLE_TIME", "15m")
	cfg.db.maxIdleTime, err = time.ParseDuration(maxIdleTimeStr)
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}
func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func (app *application) readIDParam(r *http.Request) (int64, error) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("invalid id parameter")
	}

	return id, nil
}

// readString() helper returns a string value from the query string
func (app *application) readString(qs url.Values, key string, defaultValue string) string {
	s := qs.Get(key)

	if s == "" {
		return defaultValue
	}

	return s
}

// readInt() helper returns a int value from the query string
func (app *application) readInt(qs url.Values, key string, defaultValue int, v *validator.Validator) int {
	s := qs.Get(key)

	if s == "" {
		return defaultValue
	}

	i, err := strconv.Atoi(s)
	if err != nil {
		v.AddError(key, "must be an integer value")
		return defaultValue
	}

	return i
}

func updateIfNotNil[T any](dest *T, src *T) {
	if src != nil {
		*dest = *src
	}
}
