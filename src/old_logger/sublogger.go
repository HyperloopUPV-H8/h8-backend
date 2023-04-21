package old_logger

type SubLogger interface {
	Start() error
	Stop() error
	Flush() error
	Close() error
	Update(any) error
}
