package error

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"runtime"
)

type Err struct {
	errorType Type
	code      string
	errors    []error
	message   string
	detail    map[string]interface{}
	location  string
}

type Error interface {
	json.Marshaler
	error
	ErrorCode() string
	Hash() string
	Type() Type
	WithError(error) Error
	WithMessage(string) Error
	WithDetail(string, interface{}) Error
	WithErrorLocation() Error
}

type Type string

const (
	TypeInfo     Type = "INFO"
	TypeWarn     Type = "WARNING"
	TypeError    Type = "ERROR"
	TypeCritical Type = "CRITICAL"
)

func NewError(errType Type, code string) Error {
	return Err{
		errorType: errType,
		code:      code,
		detail:    make(map[string]interface{}),
	}
}

func (e Err) ErrorCode() string {
	return e.code
}

func (e Err) Hash() string {
	h := sha256.New()
	h.Write([]byte(fmt.Sprint(e.errorType, e.code)))
	return hex.EncodeToString(h.Sum(nil))
}

func (e Err) Type() Type {
	return e.errorType
}

func (e Err) WithError(err error) Error {
	errCast, ok := err.(Err)
	if ok {
		e.errors = append(e.errors, errCast)
	} else {
		e.errors = append(e.errors, Err{
			code:      "NATIVE_ERROR",
			errorType: TypeError,
			message:   err.Error(),
		})
	}
	return e
}

func (e Err) WithMessage(message string) Error {
	e.message = message
	return e
}

func (e Err) WithDetail(detailName string, detailData interface{}) Error {
	e.detail[detailName] = detailData
	return e
}

func (e Err) WithErrorLocation() Error {
	_, f, l, _ := runtime.Caller(1)
	e.location = fmt.Sprintf("%s:%d", f, l)
	return e
}

func (e Err) Error() string {
	if e.message != "" {
		return fmt.Sprintf("[%s] %s : %s", e.errorType, e.code, e.message)
	}
	return fmt.Sprintf("[%s] %s", e.errorType, e.code)
}

func (e Err) MarshalJSON() ([]byte, error) {
	err := struct {
		Code     string                 `json:"code"`
		Type     Type                   `json:"type"`
		Errors   []error                `json:"errors,omitempty"`
		Message  string                 `json:"message,omitempty"`
		Detail   map[string]interface{} `json:"detail,omitempty"`
		Location string                 `json:"location,omitempty"`
	}{
		Code:     e.code,
		Type:     e.errorType,
		Errors:   e.errors,
		Message:  e.message,
		Detail:   e.detail,
		Location: e.location,
	}
	return json.Marshal(err)
}

// ErrorLines will return first error found from provided errors on parameter
func ErrorLines(errs ...error) error {
	for _, err := range errs {
		if err != nil {
			return err
		}
	}
	return nil
}
