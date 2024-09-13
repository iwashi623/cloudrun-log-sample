package main

import (
	"cloudrun-log-sample/mylog"
	"log/slog"
	"os"

	"github.com/labstack/echo/v4"
)

func initLogger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		slogHandler := mylog.NewHandler(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{ReplaceAttr: mylog.GoogleMessageReplacer}), mylog.SourceOption{Enabled: true, KeyName: mylog.GoogleSourceKeyName})
		slog.SetDefault(slog.New(slogHandler))

		return next(c)
	}
}

func defaultLogFunc(projectID string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			mylog.InfoContext(c, "defaultLogFunc", mylog.Args{"project_id": projectID})
			mylog.WithTrace(c, projectID)
			mylog.WithValue(c, mylog.Path, c.Path())
			mylog.WithValue(c, mylog.Method, c.Request().Method)
			mylog.WithValue(c, mylog.Query, c.QueryString())
			mylog.InfoContext(c, "Received Request")
			return next(c)
		}
	}
}
