package application

import (
	"fmt"

	"github.com/HyperloopUPV-H8/Backend-H8/message_transfer/domain"
)

type id = uint16

type MessageJSON struct {
	Id     id                `json:"id"`
	Values map[string]string `json:"values"`
}

func NewMessageJSON(message domain.Message) MessageJSON {
	return MessageJSON{
		Id:     message.ID(),
		Values: getValues(message.Values()),
	}
}

func getValues(values map[string]any) map[string]string {
	result := make(map[string]string, len(values))
	for name, value := range values {
		result[name] = fmt.Sprintf("%v", value)
	}
	return result
}
