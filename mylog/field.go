package mylog

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// echo.Contextが実体となる
// 不要なメソッドが使われないように、mylog側のInterfaceを使用する
type myContext interface {
	Request() *http.Request
	SetRequest(r *http.Request)
}

// プリミティブ型は適当な値を入れられてしまうのでstructを使う
type mylogField struct {
	keyName string
}

// stringer interface implementation
func (f *mylogField) String() string {
	return f.keyName
}

// 検索性を保つため、フィールドを制限している
var (
	// Basic fields
	Method = mylogField{keyName: "method"}
	Path   = mylogField{keyName: "path"}
	Query  = mylogField{keyName: "query"}

	// Request fields
	UserID = mylogField{keyName: "user_id"}
)

const (
	GoogleSourceKeyName = "logging.googleapis.com/sourceLocation"
	GoogleTraceKeyName  = "logging.googleapis.com/trace"
)

// GoogleMessageReplacer Google Cloud上でのログ出力に合わせて、msgをmessageに変換する置換ルール
func GoogleMessageReplacer(groups []string, a slog.Attr) slog.Attr {
	if a.Key == slog.MessageKey {
		a.Key = "message"
	}
	return a
}

type Args map[string]any

// InfoContext INFOレベルのログを出力する
func InfoContext(dealCtx myContext, msg string, args ...Args) {
	ctx := dealCtx.Request().Context()
	ctx = withValue(ctx, "severity", "INFO")
	for _, arg := range args {
		for k, v := range arg {
			ctx = withValue(ctx, k, v)
		}
	}
	slog.InfoContext(ctx, msg)
}

// WarnContext WARNINGレベルのログを出力する
func WarnContext(dealCtx myContext, msg string, args ...Args) {
	ctx := dealCtx.Request().Context()
	ctx = withValue(ctx, "severity", "WARNING")
	for _, arg := range args {
		for k, v := range arg {
			ctx = withValue(ctx, k, v)
		}
	}
	slog.WarnContext(ctx, msg)
}

// ErrorContext ERRORレベルのログを出力する
func ErrorContext(dealCtx myContext, msg string, err error, args ...Args) {
	ctx := dealCtx.Request().Context()
	ctx = withValue(ctx, "severity", "ERROR")
	for _, arg := range args {
		for k, v := range arg {
			ctx = withValue(ctx, k, v)
		}
	}

	if err == nil {
		err = errors.New("error occurred")
	}

	ctx = withValue(ctx, "stack_trace", errors.WithStack(err))
	slog.ErrorContext(ctx, fmt.Sprintf("%s err=%s", msg, err.Error()))
}

// WithValue 値をログ出力用のcontextにセットする
func WithValue(dealCtx myContext, key mylogField, value any) {
	ctx := dealCtx.Request().Context()
	ctx = withValue(ctx, key.String(), value)
	dealCtx.SetRequest(dealCtx.Request().WithContext(ctx))
}

// WithTrace traceIDをログ出力用のcontextにセットする
func WithTrace(dealCtx myContext, projectID string) {
	traceID := ""
	if traceID = getTraceID(dealCtx.Request()); traceID == "" {
		traceID = strings.Replace(uuid.New().String(), "-", "", -1) // googleのtraceIDは32文字の16進数ハイフンなしの文字列
		WarnContext(dealCtx, "traceID not found in header, generated new traceID")
	}

	trace := fmt.Sprintf("projects/%s/traces/%s", projectID, traceID)
	ctx := dealCtx.Request().Context()
	ctx = withValue(ctx, GoogleTraceKeyName, trace)
	dealCtx.SetRequest(dealCtx.Request().WithContext(ctx))
}

// headerからtraceIDを抽出する
// https://moritomo7315.hatenablog.com/entry/go-traceid-logging
func getTraceID(r *http.Request) string {
	traceHeader := r.Header.Get("X-Cloud-Trace-Context")
	traceParts := strings.Split(traceHeader, "/")
	traceID := ""
	if len(traceParts) > 0 {
		traceID = traceParts[0]
	}
	return traceID
}
