package dto

type LogValue struct {
	name string
	data string
}

func NewLogValue(name string, data string) LogValue {
	return LogValue{
		name: name,
		data: data,
	}
}

func (value LogValue) Name() string {
	return value.name
}

func (value LogValue) Data() string {
	return value.data
}
