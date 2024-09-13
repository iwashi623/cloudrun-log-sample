package main

import (
	"cloudrun-log-sample/mylog"
	"fmt"
	"log/slog"
	"math/rand"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	projectID := os.Getenv("PROJECT_ID")
	fmt.Println("projectID:", projectID)

	sampleGroup := e.Group("/simple")
	sampleGroup.GET("/:user_id", simpleUserHandler)
	sampleGroup.GET("/:user_id/with_error", simpleUserWithErrorHandler)

	slogGrop := e.Group("/slog")
	slogGrop.Use(slogSetUp)

	e.GET("/random", func(c echo.Context) error {
		fmt.Println("start halfHandler")
		// 1か0をランダムで返す
		if oneInFive() {
			return c.String(http.StatusInternalServerError, "エラーが発生しました")
		}

		return c.String(http.StatusOK, "成功しました")
	})

	e.Logger.Fatal(e.Start(":9090"))
}

func simpleUserHandler(c echo.Context) error {
	userID := c.Param("user_id")
	fmt.Println("simpleHandler user_id:", userID)
	return c.String(http.StatusOK, "simpleHandler OK, user_id: "+userID)
}

func simpleUserWithErrorHandler(c echo.Context) error {
	userID := c.Param("user_id")
	if err := hoge(); err != nil {
		fmt.Println("simpleHandlerWithErrorでエラーが発生しました user_id:", userID)
		return c.String(http.StatusInternalServerError, "simpleHandlerWithError Error")
	}
	return c.String(http.StatusOK, "simpleHandlerWithError OK")
}

func slogSetUp(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		slogHandler := mylog.NewHandler(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{ReplaceAttr: mylog.GoogleMessageReplacer}), mylog.SourceOption{Enabled: true, KeyName: mylog.GoogleSourceKeyName})
		slog.SetDefault(slog.New(slogHandler))
		return next(c)
	}
}

func hoge() error {
	return fmt.Errorf("error")
}

func oneInFive() bool {
	return rand.Intn(5) == 0
}
