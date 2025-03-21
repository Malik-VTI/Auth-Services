package middlewares

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
)

func LoggingMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		recorder := &echoResponseRecorder{ResponseWriter: c.Response().Writer}
		c.Response().Writer = recorder
		err := next(c)

		logger := c.Get("logger").(*slog.Logger)
		logger = logger.With(slog.Int("http_status", recorder.status))

		c.Set("logger", logger)

		return err
	}
}

type echoResponseRecorder struct {
	http.ResponseWriter
	status int
}

func (r *echoResponseRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}
