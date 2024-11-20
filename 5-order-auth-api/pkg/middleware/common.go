package middleware

import "net/http"

type WrapperWriter struct {
	http.ResponseWriter
	StatusCode    int
	headerWritten bool
}

func (w *WrapperWriter) WriteHeader(statusCode int) {
	if !w.headerWritten {
		w.ResponseWriter.WriteHeader(statusCode)
		w.StatusCode = statusCode
		w.headerWritten = true
	}
}
