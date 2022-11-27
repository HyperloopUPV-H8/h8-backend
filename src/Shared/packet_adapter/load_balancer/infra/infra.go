package infra

import "container/ring"

type LoadBalancer struct {
	buf *ring.Ring
}

func Init(destinations []chan<- []byte) *LoadBalancer {
	buf := ring.New(len(destinations))
	for _, dest := range destinations {
		buf.Value = dest
		buf = buf.Next()
	}

	return &LoadBalancer{
		buf: buf,
	}
}

func (balancer *LoadBalancer) Next(data []byte) {
	balancer.buf.Value.(chan<- []byte) <- data
	balancer.buf = balancer.buf.Next()
}
