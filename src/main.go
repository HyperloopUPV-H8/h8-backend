package main

import (
	"fmt"

	"github.com/HyperloopUPV-H8/Backend-H8/Shared/PacketAdapter/infra"
)

func main() {
	tc := infra.New([]string{})
	i := 0
	for {
		i++
		fmt.Println(tc.ReceiveData(), "\n-", i, "-")
	}
}
