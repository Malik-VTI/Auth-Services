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

func Auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		jsonHandler := slog.NewJSONHandler(os.Stderr, nil)
		myslog := slog.New(jsonHandler)

		cookie, err := c.Cookie("Authorization")
		if err != nil {
			if err == http.ErrNoCookie {
				myslog.Error("Token not found")
				return echo.NewHTTPError(http.StatusUnauthorized, "Token not found")
			}
			return err
		}
		tokenString := cookie.Value

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				myslog.Error("unexpected signing method: %v", t.Header["alg"])
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte("secret"), nil
		})

		if err != nil {
			myslog.Error("Invalid token")
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			if float64(time.Now().Unix()) > claims["exp"].(float64) {
				myslog.Error("Token is expired")
				return c.JSON(http.StatusUnauthorized, "Error: Token is expired")
			}

			var user model.User
			config.DB.First(&user, claims["sub"])

			if user.ID == 0 {
				myslog.Error("User not authorized")
				return c.JSON(http.StatusUnauthorized, "Ypu are not authorized, please login")
			}
			c.Set("user", user)
			return next(c)
		}
		return err
	}
}
