package handler

import (
	"database/sql"

	"github.com/labstack/echo/v4"

	echoSwagger "github.com/swaggo/echo-swagger"
)

func ApplyHandler(e *echo.Echo, db *sql.DB) {
	// 회원가입 API
	e.POST("/api/sign-up", signUp(db))

	// 로그인 API(현재는 테스트용)
	e.POST("/api/sign-in", signIn(db))

	// 애플 로그인 API(현재는 테스트용)
	e.POST("/api/apple-sign-in", appleSignIn(db))

	e.GET("/swagger/*", echoSwagger.WrapHandler)
}
