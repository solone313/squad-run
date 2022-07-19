package handler

import (
	"database/sql"
	"net/http"

	"github.com/pkg/errors"
	"github.com/solone313/squad-run/helper"
	"github.com/solone313/squad-run/models"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"

	"github.com/labstack/echo/v4"
)

// @Summary 회원가입 API
// @Description email, password를 받아 가입합니다.
// @Accept json
// @Produce json
// @Param request body SignInUserRequest true "유저 정보"
// @Success 200 {object} message
// @Failure 401 {object} message
// @Failure 400 {object} message
// @Failure 500 {object} message
// @Router /sign-up [post]
func signUp(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		req := new(SignInUserRequest)

		if err := c.Bind(req); err != nil {
			return c.JSON(http.StatusBadRequest, message{
				"bad request",
			})
		}

		user, err := models.Users(
			models.UserWhere.Email.EQ(req.Email),
		).One(ctx, db)
		if err != nil && errors.Cause(err) != sql.ErrNoRows {
			return echo.ErrInternalServerError
		}

		// 이미 이메일이 존재할 경우의 처리
		if user != nil {
			return c.JSON(http.StatusBadRequest, message{
				"existing email",
			})
		}

		// 비밀번호를 bycrypt 라이브러리로 해싱 처리
		hashpw, err := helper.HashPassword(*req.Password)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, message{
				err.Error(),
			})
		}
		req.Password = &hashpw

		// 위의 두단계에서 err가 nil일 경우 DB에 유저를 생성
		if err := (&models.User{
			Email:    user.Email,
			Password: null.StringFrom(*req.Password),
		}).Insert(ctx, db, boil.Infer()); err != nil {
			return c.JSON(http.StatusInternalServerError, message{
				err.Error(),
			})
		}

		// 모든 처리가 끝난 후 200, Success 메시지를 반환
		if err := c.JSON(http.StatusOK, message{
			"Success",
		}); err != nil {
			return errors.Wrap(err, "signup")
		}
		return nil
	}
}
