package httpx

import (
    "encoding/json"
    "net/http"
)

// Bind decodes JSON request body into T with DisallowUnknownFields and maps errors to HTTP 400.
func Bind[T any](r *http.Request) (T, error) {
    var body T
    dec := json.NewDecoder(r.Body)
    dec.DisallowUnknownFields()
    if err := dec.Decode(&body); err != nil {
        return body, BadRequest("invalid_json", err)
    }
    return body, nil
}

// EndpointJSON wraps a handler that expects a JSON body of type T.
func EndpointJSON[T any](fn func(r *http.Request, body T) (int, any, error)) http.HandlerFunc {
    return Endpoint(func(r *http.Request) (int, any, error) {
        body, err := Bind[T](r)
        if err != nil {
            return http.StatusBadRequest, nil, err
        }
        return fn(r, body)
    })
}

// Validate binds JSON and then runs validator. For 422 responses on validation errors.
func Validate[T any](r *http.Request, validate func(T) error) (T, error) {
    body, err := Bind[T](r)
    if err != nil {
        return body, err
    }
    if err := validate(body); err != nil {
        return body, Unprocessable("validation_error", err)
    }
    return body, nil
}


