package controller

import (
	"log/slog"
	"net/http"
	"os"
	"services-auth/config"
	"services-auth/model"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

var jsonHandler = slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
	Level: slog.LevelInfo,
})

var myslog = slog.New(jsonHandler)

func Register(c echo.Context) error {

	traceID := c.Request().Header.Get("X-Trace-ID")
	if traceID == "" {
		traceID = "no-trace-id"
	}

	logger := myslog.With(
		slog.String("trace_id", traceID),
		slog.String("http_method", c.Request().Method),
		slog.String("http_path", c.Path()),
	)

	var body struct {
		Username string
		Email    string
		Password string
		Company  string
	}

	if err := c.Bind(&body); err != nil {
		logger.Error("Invalid request body",
			slog.String("error", err.Error()),
			slog.Any("request_body", body),
		)
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Invalid request body",
		})
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error("Failed to generate password",
			slog.String("error", err.Error()),
		)
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Failed to generate password",
		})
	}

	newUser := model.User{
		Username: body.Username,
		Email:    body.Email,
		Password: string(hash),
		Company:  model.Company{Name: body.Company},
	}
	result := config.DB.Create(&newUser)
	if result.Error != nil {
		logger.Error("Failed to create user",
			slog.String("error", result.Error.Error()),
			slog.Any("user_data", map[string]interface{}{
				"username": body.Username,
				"email":    body.Email,
				"company":  body.Company,
			}),
		)
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Failed to create user",
		})
	}

	logger.Info("User Created",
		slog.Any("user_data", map[string]interface{}{
			"username": body.Username,
			"email":    body.Email,
			"company":  body.Company,
		}),
	)
	return c.JSON(http.StatusCreated, map[string]string{
		"message": "User created",
	})
}

func Login(c echo.Context) error {

	traceID := c.Request().Header.Get("X-Trace-ID")
	if traceID == "" {
		traceID = "no-trace-id"
	}

	logger := myslog.With(
		slog.String("trace_id", traceID),
		slog.String("http_method", c.Request().Method),
		slog.String("http_path", c.Path()),
	)

	var body struct {
		Account  string
		Password string
	}

	if err := c.Bind(&body); err != nil {
		logger.Error("Invalid request body",
			slog.String("error", err.Error()),
			slog.Any("request_body", body),
		)
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	var user model.User
	config.DB.First(&user, "username = ? OR email = ?", body.Account, body.Account)

	if user.ID == 0 {
		logger.Error("User not found",
			slog.String("account", body.Account),
		)
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "User not found",
		})
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		logger.Error("Invalid password",
			slog.String("error", err.Error()),
		)
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid password",
		})
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":      user.ID,
		"exp":      time.Now().Add(time.Hour * 24 * 30).Unix(),
		"username": user.Username,
		"email":    user.Email,
	})

	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		logger.Error("Failed to generate token",
			slog.String("error", err.Error()),
		)
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Failed to generate token",
		})
	}

	cookie := new(http.Cookie)
	cookie.Name = "Authorization"
	cookie.Value = tokenString

	cookie.SameSite = http.SameSiteLaxMode
	c.SetCookie(cookie)

	logger.Info("Login success",
		slog.String("username", user.Username),
	)
	return c.JSON(http.StatusOK, map[string]string{
		"message": "Login success",
	})
}

func Home(c echo.Context) error {
	traceID := c.Request().Header.Get("X-Trace-ID")
	if traceID == "" {
		traceID = "no-trace-id"
	}

	logger := myslog.With(
		slog.String("trace_id", traceID),
		slog.String("http_method", c.Request().Method),
		slog.String("http_path", c.Path()),
	)

	logger.Info("Home accessed")
	return c.JSON(http.StatusOK, map[string]string{
		"message": "Home",
	})

}
