package httpx

import (
    "encoding/json"
    "errors"
    "net/http"
)

// Endpoint is a tiny helper to reduce boilerplate in HTTP handlers.
// It handles setting content-type, status code, JSON encoding and basic error mapping.
func Endpoint(fn func(r *http.Request) (int, any, error)) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        status, data, err := fn(r)
        if err != nil {
            writeError(w, err)
            return
        }
        WriteJSON(w, status, data)
    }
}

func writeError(w http.ResponseWriter, err error) {
    var he *HTTPError
    if errors.As(err, &he) {
        WriteJSON(w, he.StatusCode, map[string]any{
            "error":   he.Message,
            "details": he.Details,
        })
        return
    }
    // Fallback generic error
    WriteJSON(w, http.StatusInternalServerError, map[string]any{"error": http.StatusText(http.StatusInternalServerError)})
}

// DecodeJSON decodes JSON body into dst and returns 400-compatible error on failure.
func DecodeJSON(r *http.Request, dst any) error {
    dec := json.NewDecoder(r.Body)
    dec.DisallowUnknownFields()
    if err := dec.Decode(dst); err != nil {
        return BadRequest("invalid_json", err)
    }
    return nil
}


