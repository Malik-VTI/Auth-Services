package middlewares

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"services-auth/config"
	"services-auth/model"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

var jsonHandler = slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
	Level: slog.LevelInfo,
})

var myslog = slog.New(jsonHandler)

func Auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		traceID := c.Request().Header.Get("X-Trace-ID")
		if traceID == "" {
			traceID = "no-trace-id"
		}

		logger := myslog.With(
			slog.String("trace_id", traceID),
			slog.String("http_method", c.Request().Method),
			slog.String("http_path", c.Path()),
		)

		cookie, err := c.Cookie("Authorization")
		if err != nil {
			if err == http.ErrNoCookie {
				logger.Error("Token not found",
					slog.String("error", err.Error()),
				)
				return echo.NewHTTPError(http.StatusUnauthorized, "Token not found")
			}
			return err
		}
		tokenString := cookie.Value

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				logger.Error("unexpected signing method: %v",
					t.Header["alg"],
					slog.String("error", err.Error()),
				)
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte("secret"), nil
		})

		if err != nil {
			logger.Error("Invalid token",
				slog.String("error", err.Error()),
			)
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			if float64(time.Now().Unix()) > claims["exp"].(float64) {
				logger.Error("Token is expired",
					slog.String("error", err.Error()),
				)
				return c.JSON(http.StatusUnauthorized, "Error: Token is expired")
			}

			var user model.User
			config.DB.First(&user, claims["sub"])

			if user.ID == 0 {
				logger.Error("User not authorized",
					slog.String("error", err.Error()),
				)
				return c.JSON(http.StatusUnauthorized, "Ypu are not authorized, please login")
			}
			logger.Info("User authorized",
				slog.String("username", user.Username),
			)
			c.Set("user", user)
			return next(c)
		}
		return err
	}
}
