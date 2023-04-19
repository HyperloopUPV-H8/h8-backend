package message

type Config struct {
	FaultId   uint16 `toml:"fault_id"`
	WarningId uint16 `toml:"warning_id"`
	BlcuAckId uint16 `toml:"blcu_ack_id"`
}
