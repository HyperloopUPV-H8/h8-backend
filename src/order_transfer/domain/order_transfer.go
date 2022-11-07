package domain

type OrderTransfer struct {
	OrderWebAdapterChannel chan OrderWebAdapter
}

func New() OrderTransfer {
	return OrderTransfer{
		OrderWebAdapterChannel: make(chan OrderWebAdapter),
	}
}

func (orderTransfer OrderTransfer) Invoke(sendOrder func(order OrderWebAdapter)) {
	go func() {
		for orderWebAdapter := range orderTransfer.OrderWebAdapterChannel {
			sendOrder(orderWebAdapter)
		}
	}()
}
