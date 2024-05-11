package wslog_test

import (
	"context"

	"github.com/mpuzanov/wslog"
)

func Example() {
	wslog.RemoveTime = true
	l := wslog.NewEnv("local")
	l.Info("hello")

	// Output:
	// level=INFO msg=hello
}

func ExampleAppendCtx() {
	wslog.RemoveTime = true
	l := wslog.New()
	ctx := wslog.AppendCtx(context.Background(), wslog.String("userID", "1"))
	l.InfoContext(ctx, "example1")
	ctx = wslog.AppendCtx(ctx, wslog.String("userID", "2"))
	l.InfoContext(ctx, "example2")
	ctx = wslog.AppendCtx(ctx, wslog.String("metod", "GET"))
	l.InfoContext(ctx, "example3")

	// Output:
	// level=INFO msg=example1 userID=1
	// level=INFO msg=example2 userID=2
	// level=INFO msg=example3 userID=2 metod=GET
}
