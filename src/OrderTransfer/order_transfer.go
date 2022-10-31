package ordertransfer

import (
	"log"
	"strconv"

	"github.com/HyperloopUPV-H8/Backend-H8/OrderTransfer/domain"
)

type OrderTransfer struct {
	OrderWebAdapterChannel chan OrderWebAdapter
}

func New() OrderTransfer {
	return OrderTransfer{
		OrderWebAdapterChannel: make(chan OrderWebAdapter),
	}
}

func (orderTransfer OrderTransfer) Invoke(sendOrder func(order domain.Order)) {
	go func() {
		for orderWebAdapter := range orderTransfer.OrderWebAdapterChannel {
			order := getOrder(orderWebAdapter)
			sendOrder(order)
		}
	}()
}

func getOrder(orderWA OrderWebAdapter) domain.Order {
	id, err := strconv.Atoi(orderWA.Id)

	if err != nil {
		log.Fatal("Error parsing float")
	}

	return domain.Order{
		Id:     uint16(id),
		Fields: getFields(orderWA.fields),
	}
}

func getFields(fieldsWA map[string]string) map[string]float64 {
	fields := make(map[string]float64, len(fieldsWA))
	for name, value := range fieldsWA {
		number, err := strconv.ParseFloat(value, 64)

		if err != nil {
			log.Fatal("Error parsing float")
		}

		fields[name] = number
	}

	return fields
}
