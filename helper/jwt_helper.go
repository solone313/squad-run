package helper

import (
	"database/sql"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/solone313/squad-run/models"
	"github.com/volatiletech/null/v8"
)

func CreateJWT(Email string) (string, error) {
	mySigningKey := []byte(os.Getenv("SECRET_KEY"))

	aToken := jwt.New(jwt.SigningMethodHS256)
	claims := aToken.Claims.(jwt.MapClaims)
	claims["Email"] = Email
	claims["exp"] = time.Now().Add(time.Hour * 720 * 3).Unix()

	tk, err := aToken.SignedString(mySigningKey)
	if err != nil {
		return "", errors.Wrap(err, "aToken.SignedString")
	}
	return tk, nil
}

func ValidateJWT(c echo.Context, db *sql.DB) (*models.User, error) {
	header := c.Request().Header
	authv := header.Get("Authorization")

	if authv == "" {
		return nil, errors.New("no authorization")
	}
	// Get bearer token
	if !strings.HasPrefix(strings.ToLower(authv), "bearer") {
		return nil, errors.New("invalid bearer token")
	}

	values := strings.Split(authv, " ")
	if len(values) < 2 {
		return nil, errors.New("no bearer token")
	}

	token := values[1]
	user, err := models.Users(
		models.UserWhere.AccessToken.EQ(null.StringFrom(token)),
	).One(c.Request().Context(), db)
	if err != nil && errors.Cause(err) != sql.ErrNoRows {
		return nil, err
	}
	// 존재하지않는 아이디일 경우
	if user == nil {
		return nil, echo.ErrBadRequest
	}
	if user.ExpiredAt.Time.Before(time.Now()) {
		return nil, echo.ErrUnauthorized
	}

	return user, nil
}
