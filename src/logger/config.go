package logger

type Config struct {
	BasePath        string       `toml:"base_path"`
	DataFileName    string       `toml:"data_name"`
	MessageFileName string       `toml:"message_name"`
	OrderFileName   string       `toml:"order_name"`
	FlushInterval   string       `toml:"flush_interval"`
	Topics          LoggerTopics `toml:"topics"`
}

type LoggerTopics struct {
	Enable string `toml:"enable"`
	State  string `toml:"state"`
}
