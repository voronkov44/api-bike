package middleware

import "net/http"

type WrapperWriter struct {
	http.ResponseWriter
	StatusCode int
	userEmail  string
}

func (w *WrapperWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.StatusCode = statusCode
}

func (w *WrapperWriter) Write(data []byte) (int, error) {
	if w.StatusCode == 0 {
		w.StatusCode = 200
	}
	return w.ResponseWriter.Write(data)
}

// SetEmail записывает email пользователя в wrapper (не в response headers!)
func (w *WrapperWriter) SetEmail(email string) {
	w.userEmail = email
}

// Email возвращает записанный email (или пустую строку)
func (w *WrapperWriter) Email() string {
	if w.userEmail == "" {
		return "Unauthorized"
	}
	return w.userEmail
}
