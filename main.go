package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	e.GET("/", hello)
	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})
	e.GET("/fmt", fmtHandler)
	e.Logger.Fatal(e.Start(":9090"))
}

func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

func fmtHandler(c echo.Context) error {
	fmt.Println("fmtHandler")
	return c.String(http.StatusOK, "fmt")
}
