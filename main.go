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
	e.GET("/", hello)
	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})
	e.GET("/posts", postsHandler)
	e.GET("/posts/:post_id", postHandler)
	e.GET("/fmt", fmtHandler)

	sampleGroup := e.Group("/simple")
	sampleGroup.GET("/:user_id", simpleUserHandler)
	sampleGroup.GET("/:user_id/with_error", simpleUserWithErrorHandler)

	e.GET("/random", func(c echo.Context) error {
		fmt.Println("start halfHandler")
		// 1か0をランダムで返す
		if oneInFive() {
			return c.String(http.StatusInternalServerError, "エラーが発生しました")
		}

		return c.String(http.StatusOK, "成功しました")
	})

	e.Use(slogSetUp)

	e.Logger.Fatal(e.Start(":9090"))
}

func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

func simpleUserHandler(c echo.Context) error {
	userID := c.Param("user_id")
	fmt.Println("simpleHandler user_id:", userID)
	return c.String(http.StatusOK, "user_id: "+userID)
}

func simpleUserWithErrorHandler(c echo.Context) error {
	userID := c.Param("user_id")
	fmt.Println("simpleHandler user_id:", userID)
	return c.String(http.StatusInternalServerError, "user_id: "+userID)
}

func fmtHandler(c echo.Context) error {
	fmt.Println("Start fmtHandler")
	// fmt1()
	fmt.Println("End fmtHandler")
	return c.String(http.StatusOK, "fmt")
}

// func fmt1() {
// 	fmt.Println("fmt1")
// 	fmt2()
// }

// func fmt2() {
// 	fmt.Println("fmt2")
// 	fmt3()
// }

// func fmt3() {
// 	fmt.Println("fmt3")
// }

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
