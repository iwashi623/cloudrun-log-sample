package main

import (
	"os"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	projectID := os.Getenv("PROJECT_ID")

	// 見づらいログを出すエンドポイント
	simpleGroup := e.Group("/simple")
	simpleGroup.GET("/:user_id", simpleUserHandler)
	simpleGroup.GET("/:user_id/with_error", simpleUserWithErrorHandler)
	simpleGroup.GET("/:user_id/multi_log", simpleUserMultilogHandler)

	slogGrop := e.Group("/slog")

	// 構造化ログを出すための初期化
	slogGrop.Use(initLogger)
	slogGrop.Use(defaultLogFunc(projectID))

	// 構造化ログをだすエンドポイント
	slogGrop.GET("/:user_id/multi_log", slogUserMultiHandler)

	e.Logger.Fatal(e.Start(":9090"))
}
