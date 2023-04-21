package logger_handler

type Logger interface {
	Start() error
	Log(loggable Loggable) error
	Stop() error
	Flush() error
	Close() error
}

type Loggable interface {
	Id() string
	Log() []string
}
