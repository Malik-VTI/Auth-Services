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

func Register(c echo.Context) error {
	jsonHandler := slog.NewJSONHandler(os.Stderr, nil)
	myslog := slog.New(jsonHandler)

	var body struct {
		Username string
		Email    string
		Password string
		Company  string
	}

	if err := c.Bind(&body); err != nil {
		myslog.Error("Invalid request body")
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Invalid request body",
		})
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		myslog.Error("Failed to generate password")
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
		myslog.Error("Failed to create user")
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Failed to create user",
		})
	}

	myslog.Info("User created")
	return c.JSON(http.StatusCreated, map[string]string{
		"message": "User created",
	})
}

func Login(c echo.Context) error {
	jsonHandler := slog.NewJSONHandler(os.Stderr, nil)
	myslog := slog.New(jsonHandler)

	var body struct {
		Account  string
		Password string
	}

	if err := c.Bind(&body); err != nil {
		myslog.Error("Invalid request body")
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	var user model.User
	config.DB.First(&user, "username = ? OR email = ?", body.Account, body.Account)

	if user.ID == 0 {
		myslog.Error("User not found")
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "User not found",
		})
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		myslog.Error("Invalid password")
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
		myslog.Error("Failed to generate token")
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Failed to generate token",
		})
	}

	cookie := new(http.Cookie)
	cookie.Name = "Authorization"
	cookie.Value = tokenString

	cookie.SameSite = http.SameSiteLaxMode
	c.SetCookie(cookie)

	myslog.Info("Login success")
	return c.JSON(http.StatusOK, map[string]string{
		"message": "Login success",
	})
}

func Home(c echo.Context) error {
	jsonHandler := slog.NewJSONHandler(os.Stderr, nil)
	myslog := slog.New(jsonHandler)

	myslog.Info("Home page accessed")
	return c.JSON(http.StatusOK, map[string]string{
		"message": "Home",
	})

}
