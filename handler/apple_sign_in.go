package handler

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/Timothylock/go-signin-with-apple/apple"
	"github.com/pkg/errors"
	"github.com/solone313/squad-run/helper"
	"github.com/solone313/squad-run/models"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"

	"github.com/labstack/echo/v4"
)

type AppleSignInRequest struct {
	Token string `json:"token" `
}

type SignInResponse struct {
	AccessToken string `json:"access_token"`
	ExpiredAt   string `json:"expired_at"`
}

type message struct {
	Message string `json:"message"`
}

// @Summary 애플 로그인 API
// @Description Token을 받아 access token을 반환합니다.
// @Accept json
// @Produce json
// @Param token body AppleSignInRequest true "애플로그인 token"
// @Success 200 {object} SignInResponse
// @Failure 400 {object} message
// @Failure 500 {object} message
// @Router /apple-sign-in [post]
func appleSignIn(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		token := new(AppleSignInRequest)
		if err := c.Bind(token); err != nil {
			return c.JSON(http.StatusBadRequest, message{"bad request"})
		}

		claim, _ := apple.GetClaims(token.Token)
		email := (*claim)["email"].(string)
		emailVerified := (*claim)["email_verified"].(string)
		if emailVerified != "true" {
			return c.JSON(http.StatusBadRequest, message{"인증되지 않은 이메일입니다."})
		}

		accessToken, err := helper.CreateJWT(email)
		if err != nil {
			return echo.ErrInternalServerError
		}

		user, err := models.Users(
			models.UserWhere.Email.EQ(email),
		).One(ctx, db)
		if err != nil && errors.Cause(err) != sql.ErrNoRows {
			return echo.ErrInternalServerError
		}

		expiredAt := time.Now().AddDate(0, 3, 0)
		// 이미 이메일이 존재할 경우의 처리
		if user != nil {
			user.AccessToken = null.StringFrom(accessToken)
			user.ExpiredAt = null.TimeFrom(expiredAt)
			if _, err := user.Update(ctx, db, boil.Infer()); err != nil {
				return c.JSON(http.StatusInternalServerError, message{"Failed insert user"})
			}
		} else {
			if err := (&models.User{
				Email:       email,
				AccessToken: null.StringFrom(accessToken),
				ExpiredAt:   null.TimeFrom(expiredAt),
			}).Insert(ctx, db, boil.Infer()); err != nil {
				return c.JSON(http.StatusInternalServerError, message{"Failed insert user"})
			}
		}

		if err := c.JSON(http.StatusOK, SignInResponse{
			AccessToken: accessToken,
			ExpiredAt:   expiredAt.Format("2006-01-02"),
		}); err != nil {
			return errors.Wrap(err, "signIn")
		}
		return nil
	}
}
