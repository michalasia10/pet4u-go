package httpx

import (
    "reflect"

    "github.com/go-playground/validator/v10"
)

var validatorInstance = validator.New()

// ValidateTags validates a struct using `validate:"..."` tags.
// On error it returns HTTP 422 with a map of field->rule.
func ValidateTags(body any) error {
    if err := validatorInstance.Struct(body); err != nil {
        if ve, ok := err.(validator.ValidationErrors); ok {
            details := map[string]string{}
            for _, e := range ve {
                field := jsonFieldName(e.StructField(), e.StructField(), body)
                details[field] = e.Tag()
            }
            return Unprocessable("validation_error", details)
        }
        return Unprocessable("validation_error", err.Error())
    }
    return nil
}

// jsonFieldName attempts to map a struct field to its json tag name.
func jsonFieldName(structField string, fallback string, body any) string {
    t := reflect.TypeOf(body)
    if t.Kind() == reflect.Pointer { t = t.Elem() }
    if t.Kind() != reflect.Struct { return fallback }
    if f, ok := t.FieldByName(structField); ok {
        tag := f.Tag.Get("json")
        if tag != "" && tag != "-" {
            // trim ",omitempty" etc.
            if comma := indexByte(tag, ','); comma > 0 { return tag[:comma] }
            return tag
        }
    }
    return fallback
}

func indexByte(s string, c byte) int {
    for i := 0; i < len(s); i++ {
        if s[i] == c { return i }
    }
    return -1
}


