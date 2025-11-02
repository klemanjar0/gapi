package response_builder

import (
	"encoding/json"
	"net/http"
)

type ResponseBuilder struct {
	w          http.ResponseWriter
	statusCode int
	data       interface{}
}

type ResponseOption func(*ResponseBuilder)

func NewResponseBuilder(w http.ResponseWriter) *ResponseBuilder {
	return &ResponseBuilder{
		w:          w,
		statusCode: http.StatusOK,
	}
}

func WithStatusCode(code int) ResponseOption {
	return func(rb *ResponseBuilder) {
		rb.statusCode = code
	}
}

func WithData(data interface{}) ResponseOption {
	return func(rb *ResponseBuilder) {
		rb.data = data
	}
}

func WithError(message string) ResponseOption {
	return func(rb *ResponseBuilder) {
		if message == "" {
			message = "Unknown error occurred"
		}
		rb.data = map[string]string{"error": message}
	}
}

func (rb *ResponseBuilder) Success(data interface{}) {
	rb.WriteResponse(WithData(data), WithStatusCode(http.StatusOK))
}

func (rb *ResponseBuilder) Created(data interface{}) {
	rb.WriteResponse(WithData(data), WithStatusCode(http.StatusCreated))
}

func (rb *ResponseBuilder) Error(message string, code int) {
	rb.WriteResponse(WithError(message), WithStatusCode(code))
}

func (rb *ResponseBuilder) WriteResponse(opts ...ResponseOption) {
	for _, opt := range opts {
		opt(rb)
	}

	if rb.statusCode == http.StatusNoContent {
		rb.w.WriteHeader(http.StatusNoContent)
		return
	}

	rb.w.Header().Set("Content-Type", "application/json")
	rb.w.WriteHeader(rb.statusCode)

	if rb.data != nil {
		json.NewEncoder(rb.w).Encode(rb.data)
	}
}
