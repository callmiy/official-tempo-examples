package log

import (
	"os"

	kitlog "github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/weaveworks/common/logging"
	"github.com/weaveworks/common/server"
)

// Logger is a shared go-kit logger.
// TODO: Change all components to take a non-global logger via their constructors.
// Prefer accepting a non-global logger as an argument.
var Logger = kitlog.NewNopLogger()

// InitLogger initialises the global gokit logger and overrides the
// default logger for the server.
func InitLogger(cfg *server.Config) {
	logger := kitlog.NewLogfmtLogger(kitlog.NewSyncWriter(os.Stderr))
	if cfg.LogFormat.String() == "json" {
		logger = kitlog.NewJSONLogger(kitlog.NewSyncWriter(os.Stderr))
	}

	// add support for level based logging
	logger = level.NewFilter(logger, LevelFilter(cfg.LogLevel.String()))

	// use UTC timestamps
	logger = kitlog.With(logger, "ts", kitlog.DefaultTimestampUTC)

	// when use util_log.Logger, skip 3 stack frames.
	Logger = kitlog.With(logger, "caller", kitlog.Caller(3))

	// cfg.Log wraps log function, skip 4 stack frames to get caller information.
	// this works in go 1.12, but doesn't work in versions earlier.
	// it will always shows the wrapper function generated by compiler
	// marked <autogenerated> in old versions.
	cfg.Log = logging.GoKit(kitlog.With(logger, "caller", kitlog.Caller(4)))
}

// TODO: remove once weaveworks/common updates to go-kit/log
// -> we can then revert to using Level.Gokit
func LevelFilter(l string) level.Option {
	switch l {
	case "debug":
		return level.AllowDebug()
	case "info":
		return level.AllowInfo()
	case "warn":
		return level.AllowWarn()
	case "error":
		return level.AllowError()
	default:
		return level.AllowAll()
	}
}
