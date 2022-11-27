package tcp

type Config struct {
	LocalPort   uint16
	RemoteIPs   []string
	RemotePorts []uint16
	Snaplen     int32
}

func DefaultConfig() *Config {
	return &Config{
		LocalPort:   50000,
		RemoteIPs:   []string{"127.0.0.1"},
		RemotePorts: []uint16{50001},
		Snaplen:     ^int32(0),
	}
}
