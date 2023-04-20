package data_logger

type Config struct {
	BasePath string `toml:"base_path"`
	FileName string `toml:"file_name"`
}
