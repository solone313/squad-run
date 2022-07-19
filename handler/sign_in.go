package handler

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"github.com/solone313/squad-run/helper"
	"github.com/solone313/squad-run/models"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"

	"github.com/labstack/echo/v4"
)

type SignInUserRequest struct {
	Email    string  `json:"email"`
	Password *string `json:"password"`
}

// @Summary 로그인 API
// @Description email, password를 받아 access token을 반환합니다.
// @Accept json
// @Produce json
// @Param request body SignInUserRequest true "유저 정보"
// @Success 200 {object} SignInResponse
// @Failure 401 {object} message
// @Failure 400 {object} message
// @Failure 500 {object} message
// @Router /sign-in [post]
func signIn(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		req := new(SignInUserRequest)

		if err := c.Bind(req); err != nil {
			return c.JSON(http.StatusBadRequest, message{"bad request"})
		}
		inputpw := req.Password

		user, err := models.Users(
			models.UserWhere.Email.EQ(req.Email),
		).One(ctx, db)
		if err != nil && errors.Cause(err) != sql.ErrNoRows {
			return echo.ErrInternalServerError
		}

		// 존재하지않는 아이디일 경우
		if user == nil {
			return echo.ErrBadRequest
		}

		res := helper.CheckPasswordHash(user.Password.String, *inputpw)

		// 비밀번호 검증에 실패한 경우
		if !res {
			return echo.ErrUnauthorized
		}
		// 토큰 발행
		accessToken, err := helper.CreateJWT(user.Email)
		if err != nil {
			return echo.ErrInternalServerError
		}
		expiredAt := time.Now().AddDate(0, 3, 0)
		user.AccessToken = null.StringFrom(accessToken)
		user.ExpiredAt = null.TimeFrom(expiredAt)
		if _, err := user.Update(ctx, db, boil.Infer()); err != nil {
			return c.JSON(http.StatusInternalServerError, message{"Failed insert user"})
		}

		if err := c.JSON(http.StatusOK, SignInResponse{
			AccessToken: accessToken,
			ExpiredAt:   user.ExpiredAt.Time.Format("2006-01-02"),
		}); err != nil {
			return errors.Wrap(err, "signIn")
		}
		return nil
	}
}
