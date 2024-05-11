package wslog

import (
	"context"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

const fileNameLogDefault = "log.txt"

type (
	// Logger алиасы типов
	Logger = slog.Logger
	// Attr ...
	Attr = slog.Attr
	// Level ...
	Level = slog.Level
	// Handler ...
	Handler = slog.Handler
	// HandlerOptions ...
	HandlerOptions = slog.HandlerOptions
	// Value ...
	Value = slog.Value
	// LogValuer ...
	LogValuer = slog.LogValuer
)

type options struct {
	Level      slog.Level
	AddSource  bool
	IsJSON     bool
	Writer     io.Writer
	FileLog    string
	OnlyFile   bool
	SetDefault bool
}

// LoggerOption ...
type LoggerOption func(options *options)

var (
	// SetDefault ...
	SetDefault = slog.SetDefault

	// String алиасы типов
	String = slog.String
	// Bool ...
	Bool = slog.Bool
	// Float64 ...
	Float64 = slog.Float64
	// Any ...
	Any = slog.Any
	// Duration ...
	Duration = slog.Duration
	// Int ...
	Int = slog.Int
	// Int64 ...
	Int64 = slog.Int64

	// GroupValue ...
	GroupValue = slog.GroupValue
	// Group ...
	Group = slog.Group
)

const (
	defaultLevel      = slog.LevelInfo
	defaultSetDefault = true
)

var (
	handler        slog.Handler
	handlerOptions *slog.HandlerOptions
	config         options
	// RemoveTime убрать time из логов
	RemoveTime bool
)

var replaceAttr = func(groups []string, a slog.Attr) slog.Attr {
	//fmt.Printf("replace: Key:%v, Value:%v,  type: %T\n", a.Key, a.Value, a.Value.Any())

	if a.Key == slog.TimeKey {
		if RemoveTime {
			return slog.Attr{}
		}
		t, ok := a.Value.Any().(time.Time)
		if ok {
			a.Value = slog.StringValue(t.Format("2006-01-02 15:04:05.000"))
		}
	}

	// Remove the directory from the source's filename.
	if a.Key == slog.SourceKey {
		source := a.Value.Any().(*slog.Source)
		source.File = filepath.Base(source.File)
		source.Function = strings.TrimLeft(filepath.Ext(source.Function), ".")
	}
	// Remove empty msg
	if a.Key == slog.MessageKey && a.Value.String() == "" {
		return slog.Attr{}
	}
	return a
}

// New create logger
func New(opts ...LoggerOption) *Logger {

	config = options{
		Level:      defaultLevel,
		Writer:     os.Stdout,
		SetDefault: defaultSetDefault,
		AddSource:  false,
	}

	for _, opt := range opts {
		opt(&config)
	}

	if config.FileLog != "" {
		fw, _ := os.OpenFile(config.FileLog, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if config.OnlyFile {
			config.Writer = fw
		} else {
			config.Writer = io.MultiWriter(config.Writer, fw)
		}
	}
	log := newLogger()
	return log
}

func newLogger() *Logger {

	handlerOptions = &slog.HandlerOptions{
		Level:       config.Level,
		AddSource:   config.AddSource,
		ReplaceAttr: replaceAttr,
	}

	if config.IsJSON {
		handler = slog.NewJSONHandler(config.Writer, handlerOptions)
	} else {
		handler = slog.NewTextHandler(config.Writer, handlerOptions)
	}

	h := &ContextHandler{handler}
	logger := slog.New(h)

	if config.SetDefault {
		slog.SetDefault(logger)
	}

	return logger
}

// WithFileLog ...
func WithFileLog(fileName string) LoggerOption {
	return func(options *options) {
		if fileName != "" {
			options.FileLog = fileName
		}
	}
}

// WithOnlyFile ...
func WithOnlyFile(onlyFile bool) LoggerOption {
	return func(options *options) {
		options.OnlyFile = onlyFile
	}
}

// WithIsJSON ...
func WithIsJSON(IsJSON bool) LoggerOption {
	return func(options *options) {
		options.IsJSON = IsJSON
	}
}

// WithLevel ...
func WithLevel(level string) LoggerOption {
	return func(options *options) {
		var l slog.Level
		if err := l.UnmarshalText([]byte(level)); err != nil {
			l = slog.LevelInfo
		}
		options.Level = l
	}
}

// WithAddSource logger option sets the add source option, which will add source file and line number to the log record.
func WithAddSource(addSource bool) LoggerOption {
	return func(o *options) {
		o.AddSource = addSource
	}
}

// WithWriter ...
func WithWriter(w io.Writer) LoggerOption {
	return func(o *options) {
		o.Writer = w
	}
}

// NewEnv создание логера на основе Env (local, dev, prod)
func NewEnv(env string) *Logger {

	config = options{
		Level:      defaultLevel,
		Writer:     os.Stdout,
		SetDefault: defaultSetDefault,
		IsJSON:     true,
		AddSource:  true,
	}

	switch env {

	case envLocal:
		config.Level = slog.LevelDebug
		config.IsJSON = false
		config.AddSource = false

	case envDev:
		config.Level = slog.LevelDebug

	case envProd:
		config.Level = slog.LevelInfo

		// добавим запись логов в файл
		file := &lumberjack.Logger{
			Filename:   fileNameLogDefault, // Имя файла
			MaxSize:    10,                 // Размер в МБ до ротации файла
			MaxBackups: 5,                  // Максимальное количество файлов, сохраненных до перезаписи
			MaxAge:     30,                 // Максимальное количество дней для хранения файлов
			Compress:   true,               // Следует ли сжимать файлы логов с помощью gzip
		}

		config.Writer = io.MultiWriter(file, os.Stdout) //, os.Stderr

	default:
		config.Level = defaultLevel
	}

	log := newLogger()
	return log
}

//=======================================================

// ErrAttr helper func
//
// logger.Error("user msg error", wslog.ErrAttr(err))
func ErrAttr(err error) Attr {
	return slog.String("error", err.Error())
}

// Default ...
func Default() *Logger {
	return slog.Default()
}

// SetLogLevel ...
func SetLogLevel(level string) slog.Level {
	var l slog.Level
	if err := l.UnmarshalText([]byte(level)); err != nil {
		l = slog.LevelInfo
	}
	config.Level = l
	newLogger()
	return l
}

// GetLogLevel ...
func GetLogLevel() slog.Level {
	return config.Level
}

//=======================================================

type ctxLogger struct{}

// ContextWithLogger adds logger to context
func ContextWithLogger(ctx context.Context, l *slog.Logger) context.Context {
	return context.WithValue(ctx, ctxLogger{}, l)
}

// loggerFromContext returns logger from context
func loggerFromContext(ctx context.Context) *slog.Logger {
	if l, ok := ctx.Value(ctxLogger{}).(*slog.Logger); ok {
		return l
	}
	return slog.Default()
}

// L loggerFromContext
func L(ctx context.Context) *Logger {
	return loggerFromContext(ctx)
}

// WithAttrs returns logger with attributes.
func WithAttrs(ctx context.Context, attrs ...Attr) *Logger {
	logger := L(ctx)
	for _, attr := range attrs {
		logger = logger.With(attr)
	}

	return logger
}

// =======================================================
