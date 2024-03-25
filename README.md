# Wrapper for slog

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

or

logger := wslog.NewEnv(cfg.Env)
logger.Debug("debug", wslog.AnyAttr("cfg", cfg))

logger.Info(cfg.ServiceName,
		wslog.StrAttr("version", config.Version),
		wslog.StrAttr("time", config.Time),
		wslog.StrAttr("log_level", cfg.Log.Level),
	)
```