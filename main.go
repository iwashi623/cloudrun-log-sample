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

	simpleGroup := e.Group("/simple")
	simpleGroup.GET("/:user_id", simpleUserHandler)
	simpleGroup.GET("/:user_id/with_error", simpleUserWithErrorHandler)
	simpleGroup.GET("/:user_id/multi_log", simpleUserMultilogHandler)

	slogGrop := e.Group("/slog")
	slogGrop.Use(initLogger)
	slogGrop.Use(defaultLogFunc(projectID))
	slogGrop.GET("/:user_id/multi_log", slogUserMultiHandler)

	e.Logger.Fatal(e.Start(":9090"))
}

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

func hoge() error {
	return fmt.Errorf("error")
}

// ログに出すための適当な文字列が入ったスライス
var strs = []string{
	"hoge",
	"fuga",
	"piyo",
	"foo",
	"bar",
	"baz",
	"qux",
	"quux",
	"corge",
	"grault",
	"garply",
	"waldo",
	"fred",
	"plugh",
	"xyzzy",
	"thud",
	"hogehoge",
	"fugafuga",
	"piyopiyo",
	"foofoo",
	"barbar",
	"bazbaz",
	"quxqux",
	"quuxquux",
	"corgecorge",
	"graultgrault",
	"garplygarply",
	"waldowaldo",
	"fredfred",
	"plughplugh",
	"xyzzyxyzzy",
	"thudthud",
	"hogehogehoge",
	"fugafugafuga",
	"piyopiyopiyo",
	"foofoofoo",
	"barbarbar",
	"bazbazbaz",
	"quxquxqux",
	"quuxquuxquux",
	"corgecorgecorge",
	"graultgraultgrault",
	"garplygarplygarply",
	"waldowaldowaldo",
	"fredfredfred",
	"plughplughplugh",
	"xyzzyxyzzyxyzzy",
	"thudthudthud",
	"hogehogehogehoge",
	"fugafugafugafuga",
	"piyopiyopiyopiyo",
	"foofoofoofoo",
	"barbarbarbar",
	"bazbazbazbaz",
	"quxquxquxqux",
	"quuxquuxquuxquux",
	"corgecorgecorgecorge",
	"graultgraultgraultgrault",
}

func multiLogByFmt() error {
	count := rand.Intn(50)
	for i := 0; i < count; i++ {
		fmt.Println(strs[rand.Intn(len(strs))])
	}
	return fmt.Errorf("error")
}

func multiLogBySlog(ctx echo.Context) error {
	count := rand.Intn(50)
	for i := 0; i < count; i++ {
		mylog.InfoContext(ctx, strs[rand.Intn(len(strs))])
	}
	return fmt.Errorf("error")
}

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

func slogUserMultiHandler(c echo.Context) error {
	userID := c.Param("user_id")
	mylog.WithValue(c, mylog.UserID, userID)

	if err := multiLogBySlog(c); err != nil {
		mylog.ErrorContext(c, "slogUserMultiHandlerでエラーが発生しました", err)
		return c.String(http.StatusInternalServerError, "slogUserMultiHandler Error")
	}
	return c.String(http.StatusOK, "slogUserMultiHandler OK")
}
