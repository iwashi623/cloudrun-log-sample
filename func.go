package main

import (
	"cloudrun-log-sample/mylog"
	"fmt"
	"math/rand"

	"github.com/labstack/echo/v4"
)

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
