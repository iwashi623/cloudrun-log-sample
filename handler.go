package main

import (
	"cloudrun-log-sample/mylog"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func simpleUserHandler(c echo.Context) error {
	return c.String(http.StatusOK, "simpleHandler OK")
}

func simpleUserWithErrorHandler(c echo.Context) error {
	userID := c.Param("user_id")
	if err := hoge(); err != nil {
		fmt.Println("simpleHandlerWithErrorでエラーが発生しました user_id:", userID)
		return c.String(http.StatusInternalServerError, "simpleHandlerWithError Error")
	}
	return c.String(http.StatusOK, "simpleHandlerWithError OK")
}

func simpleUserMultilogHandler(c echo.Context) error {
	userID := c.Param("user_id")
	if err := multiLogByFmt(); err != nil {
		fmt.Println("simpleUserMultilogHandlerでエラーが発生しました user_id:", userID)
		return c.String(http.StatusInternalServerError, "simpleUserMultilogHandler Error")
	}
	return c.String(http.StatusOK, "simpleUserMultilogHandler OK")
}

func slogUserMultiHandler(c echo.Context) error {

	userID := c.Param("user_id")
	mylog.WithValue(c, mylog.UserID, userID)

	if err := multiLogBySlog(c); err != nil {
		mylog.ErrorContext(c, "slogUserMultiHandlerでエラーが発生しました", err)
		return c.String(http.StatusInternalServerError, "slogUserMultiHandler Error")
	}
	return c.String(http.StatusOK, "slogUserMultiHandler OK")
}
