package resource

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

func JSON(w http.ResponseWriter, code int, body interface{}) {
	w.Header().Set(HeaderContentType, MIMEApplicationJSON)
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(body)
}

func JSONError(w http.ResponseWriter, code int, err error) {
	res := ErrorResponse{
		Code:    code,
		Message: err.Error(),
	}
	unwrapped := errors.Unwrap(err)
	if unwrapped != nil {
		res.Error = unwrapped
	}
	JSON(w, code, res)
}

func parseParamsInt(r *http.Request, key string, defaultValue int) (int, error) {
	raw := r.URL.Query().Get(key)
	if raw == "" {
		return defaultValue, nil
	}
	val, err := strconv.Atoi(raw)
	if err != nil {
		return 0, fmt.Errorf("failed to parse integer parameter %q: %w", key, err)
	}
	return val, nil
}
