# Wrapper for slog

## Install

>go get github.com/mpuzanov/wslog

## Added methods

```golang
// NewEnv Creating a logger based on Env (local, dev, prod)  
func NewEnv(env string) *Logger

// New create logger  
func New(opts ...LoggerOption) *Logger

```

## Examples

```golang
logger := wslog.New(
		wslog.WithLevel(cfg.Log.Level),
		wslog.WithIsJSON(cfg.Log.IsJSON),
		wslog.WithFileLog(cfg.Log.File),
		wslog.WithAddSource(cfg.Log.AddSource),
	)
```

or

```golang
logger := wslog.NewEnv(cfg.Env)
logger.Debug("debug", wslog.AnyAttr("cfg", cfg))

logger.Info(cfg.ServiceName,
		wslog.StrAttr("version", config.Version),
		wslog.StrAttr("time", config.Time),
		wslog.StrAttr("log_level", cfg.Log.Level),
	)
```

adds an slog attribute to the provided context

```golang

ctx := r.Context()
reqID := uuid.New().String()
ctx = wslog.AppendCtx(ctx, slog.String("reqID", reqID))

slog.InfoContext(ctx, "request",
	slog.String("user_ip", r.RemoteAddr),
	slog.String("path", r.URL.Path),
	slog.String("method", r.Method),
)

```