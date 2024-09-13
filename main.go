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
	e.GET("/", hello)
	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})
	e.GET("/posts", postsHandler)
	e.GET("/posts/:post_id", postHandler)

	sampleGroup := e.Group("/simple")
	sampleGroup.GET("/:user_id", simpleUserHandler)
	sampleGroup.GET("/:user_id/with_error", simpleUserWithErrorHandler)

	slogGrop := e.Group("/slog")
	slogGrop.Use(slogSetUp)
	slogGrop.GET("/hello", hello)
	slogGrop.GET("/posts", postsHandler)

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

func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

func simpleUserHandler(c echo.Context) error {
	userID := c.Param("user_id")
	fmt.Println("simpleHandler user_id:", userID)
	return c.String(http.StatusOK, "simpleHandler OK, user_id: "+userID)
}

func simpleUserWithErrorHandler(c echo.Context) error {
	userID := c.Param("user_id")
	fmt.Println("simpleHandlerWithError user_id:", userID)
	return c.String(http.StatusInternalServerError, "simpleHandlerWithError Error, user_id: "+userID)
}

func postsHandler(c echo.Context) error {
	fmt.Println("postsHandler")
	return c.String(http.StatusOK, "posts")
}

func postHandler(c echo.Context) error {
	id := c.Param("post_id")
	fmt.Println("postHandler post_id:", id)
	return c.String(http.StatusOK, "post_id: "+id)
}

func slogSetUp(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		slogHandler := mylog.NewHandler(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{ReplaceAttr: mylog.GoogleMessageReplacer}), mylog.SourceOption{Enabled: true, KeyName: mylog.GoogleSourceKeyName})
		slog.SetDefault(slog.New(slogHandler))
		return next(c)
	}
}

func oneInFive() bool {
	return rand.Intn(5) == 0
}
