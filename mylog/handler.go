package mylog

import (
	"context"
	"log/slog"
	"runtime"
	"sync"
)

var _ slog.Handler = (*Handler)(nil)

type Handler struct {
	handler      slog.Handler
	sourceOption SourceOption
}

type SourceOption struct {
	Enabled bool
	KeyName string
}

// NewHandler Handlerを生成する
func NewHandler(handler slog.Handler, sourceOption SourceOption) slog.Handler {
	return Handler{
		handler:      handler,
		sourceOption: sourceOption,
	}
}

// Enabled ログ出力が有効かどうかを返す
func (h Handler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

// Handle ログ出力用のcontextにセットされた値を取得し、ログエントリに追加する
func (h Handler) Handle(ctx context.Context, record slog.Record) error {
	if v, ok := ctx.Value(fields).(*sync.Map); ok {
		v.Range(func(key, val any) bool {
			if keyString, ok := key.(string); ok {
				record.AddAttrs(slog.Any(keyString, val))
			}
			return true
		})
	}

	if h.sourceOption.Enabled {
		file, line, fn := getSource(5)
		record.AddAttrs(slog.Group(h.sourceOption.KeyName,
			slog.String("file", file),
			slog.Int("line", line),
			slog.String("function", fn),
		))
	}

	return h.handler.Handle(ctx, record)
}

// WithAttrs ログエントリに追加の属性を追加する
func (h Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return Handler{h.handler.WithAttrs(attrs), h.sourceOption}
}

// WithGroup ログエントリにグループを追加する
func (h Handler) WithGroup(name string) slog.Handler {
	return h.handler.WithGroup(name)
}

// getSource ログが呼び出されたファイル名、行番号、関数名を取得する
// hierarchyは呼び出し元の階層を指定する
func getSource(hierarchy int) (file string, line int, fn string) {
	pc, pwd, line, _ := runtime.Caller(hierarchy)
	fn = runtime.FuncForPC(pc).Name()
	return pwd, line, fn
}
