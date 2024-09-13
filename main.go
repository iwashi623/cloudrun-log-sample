package main

import (
	"cloudrun-log-sample/mylog"
	"fmt"
	"log/slog"
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

	e.Use(slogSetUp)

	e.Logger.Fatal(e.Start(":9090"))
}

func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

func fmtHandler(c echo.Context) error {
	fmt.Println("fmtHandler")
	return c.String(http.StatusOK, "fmt")
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
