# Propuesta de Logger

El logger esta compuesto por subloggers, cada uno identificado un ID ("Data", "Order", "Messages")

```go
type Loggable interface {
    Id() string
    Log() string
}

// ----------- PACKET --------------

type PacketUpdate struct {
	Metadata packet.Metadata
	HexValue []byte
	Values   map[string]packet.Value
}

type LoggablePacket struct {
    Metadata packet.Metadata
    HexValue []byte
}

func(packet LoggablePacket) Id() string {
    return packet.Metadata.Id.toString()
}

// ----------- VALUES --------------

type LoggableValue struct {
    Id string // se coge del mapa Values de PacketUpdate
    Value packet.Value
    Timestamp time.Time // la del PacketUpdate
}

func (value LoggableValue) Id() string {
    return value.Id
}

func (value LoggableValue) Log() string {
    var log string
    log += value.Timestamp + "\t"
    log += value.Value + "\n"
}

// ----------- ORDER --------------

type Order struct {
	ID     uint16           `json:"id"`
	Fields map[string]Field `json:"fields"`
    Timestamp time.Time // se genera cuando se recibe del front
}

type Field struct {
	Value     any  `json:"value"`
	IsEnabled bool `json:"isEnabled"`
}

type LoggableOrder struct {
    Id uint16
    Fields  map[string]Field
    Timestamp time.Time
}

func (order LoggableOrder) Id() string {
    return order.ID
}

func (order LoggableOrder) Log() string {
    var log string
    log += order.ID + "\t"

    for name, value of order.Fields {
        log += name + "\t" + value
    }

    log += "\n"

    return log
}

// ----------- PROTECTION --------------

type Protection struct {
	Kind      string    `json:"kind"`
	Board     string    `json:"board"`
	Value     string    `json:"value"`
	Violation Violation `json:"violation"`
	Timestamp Timestamp `json:"timestamp"`
}

type Timestamp struct {
	Counter uint16 `json:"counter"`
	Seconds uint8  `json:"seconds"`
	Minutes uint8  `json:"minutes"`
	Hours   uint8  `json:"hours"`
	Day     uint8  `json:"day"`
	Month   uint8  `json:"month"`
	Year    uint8  `json:"year"`
}

type LoggableProtection Protection

func (protection LoggableProtection) Id() string {
    return protection.Kind
}

func (protection LoggableProtection) Log() string {
    var log string

    log += protection.Board + "\t"
    log += protection.Value + "\t"
    log += LoggableViolation(protection.Violation).Log() + "\t"
    log += LoggableTimestamp(protection.Timestamp).Log() + "\n"

    return log
}

type LoggableTimestamp Timestamp

func(t LoggableTimestamp) Log() string {
    var log string

    log += fmt.Sprinf("%d/%d/%d %d:%d:%d", t.Day, t.Month, t.Year, t.Hours, t.Minutes, t.Seconds)
    log += \n

    // Alternativa
    log += t.Counter // counter es la precision restante o el timestamp al completo en micras?
}


```
