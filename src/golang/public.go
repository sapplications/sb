package golang

type Builder struct {
	Logger Logger
	items  map[string]map[string]string
}

type Generator struct {
	Logger Logger
	items  map[string]map[string]string
}

type Logger interface {
	Trace(msg string, args ...interface{})
	Debug(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	IsTrace() bool
	IsDebug() bool
	IsInfo() bool
	IsWarn() bool
	IsError() bool
}
