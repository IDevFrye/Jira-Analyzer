package middleware

import (
	"context"
	"net/http"
	"time"

	"log/slog"
)

// RequestIDKey тип для ключа контекста с ID запроса
type RequestIDKey struct{}

// LoggerConfig конфигурация для middleware логгера
type LoggerConfig struct {
	LogRequestID   bool
	LogUserAgent   bool
	LogRequestBody bool
}

// DefaultLoggerConfig возвращает конфигурацию по умолчанию
func DefaultLoggerConfig() LoggerConfig {
	return LoggerConfig{
		LogRequestID:   true,
		LogUserAgent:   true,
		LogRequestBody: false,
	}
}

// NewLoggerMiddleware создает новую middleware для логирования запросов
func NewLoggerMiddleware(log *slog.Logger, config ...LoggerConfig) func(next http.Handler) http.Handler {
	cfg := DefaultLoggerConfig()
	if len(config) > 0 {
		cfg = config[0]
	}

	return func(next http.Handler) http.Handler {
		log := log.With(
			slog.String("component", "middleware/logger"),
		)

		log.Info("logger middleware enabled")

		fn := func(w http.ResponseWriter, r *http.Request) {
			// Собираем атрибуты для логирования
			attrs := []any{
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.String("remote_addr", r.RemoteAddr),
			}

			// Добавляем опциональные атрибуты
			if cfg.LogUserAgent {
				attrs = append(attrs, slog.String("user_agent", r.UserAgent()))
			}

			if cfg.LogRequestID {
				if requestID := GetRequestID(r.Context()); requestID != "" {
					attrs = append(attrs, slog.String("request_id", requestID))
				}
			}

			// Создаем логгер с атрибутами
			entry := log.With(attrs...)

			// Создаем обертку для ResponseWriter
			ww := NewWrapResponseWriter(w)

			t1 := time.Now()
			defer func() {
				entry.Info("request completed",
					slog.Int("status", ww.Status()),
					slog.Int("bytes", ww.BytesWritten()),
					slog.String("duration", time.Since(t1).String()),
				)
			}()

			next.ServeHTTP(ww, r)
		}

		return http.HandlerFunc(fn)
	}
}

// GetRequestID извлекает ID запроса из контекста
func GetRequestID(ctx context.Context) string {
	if id, ok := ctx.Value(RequestIDKey{}).(string); ok {
		return id
	}
	return ""
}

// WrapResponseWriter обертка для http.ResponseWriter
type WrapResponseWriter struct {
	http.ResponseWriter
	status      int
	bytes       int
	wroteHeader bool
}

// NewWrapResponseWriter создает новую обертку
func NewWrapResponseWriter(w http.ResponseWriter) *WrapResponseWriter {
	return &WrapResponseWriter{ResponseWriter: w}
}

// Status возвращает HTTP статус
func (w *WrapResponseWriter) Status() int {
	return w.status
}

// BytesWritten возвращает количество байт
func (w *WrapResponseWriter) BytesWritten() int {
	return w.bytes
}

// WriteHeader перехватывает статус
func (w *WrapResponseWriter) WriteHeader(code int) {
	if !w.wroteHeader {
		w.status = code
		w.wroteHeader = true
		w.ResponseWriter.WriteHeader(code)
	}
}

// Write перехватывает запись
func (w *WrapResponseWriter) Write(b []byte) (int, error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	n, err := w.ResponseWriter.Write(b)
	w.bytes += n
	return n, err
}
