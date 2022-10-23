package PacketAdapter

import (
	"fmt"
	"net"
	"testing"
	"time"
)

// Make sure NewTransportController doesnt panic during creation
func TestNewTransportController(t *testing.T) {
	type testCase struct {
		Name  string
		Input []string
	}

	cases := []testCase{
		{
			Name:  "no addresses",
			Input: []string{},
		},
		{
			Name:  "one address",
			Input: []string{"127.0.0.1"},
		},
		{
			Name:  "invalid address",
			Input: []string{"abcd"},
		},
		{
			Name:  "multiple invalid addresses",
			Input: []string{"abcd", "efgh"},
		},
		{
			Name:  "valid and invalid addresses",
			Input: []string{"127.0.0.1", "abcd"},
		},
		{
			Name:  "duplicated addresses",
			Input: []string{"127.0.0.1", "127.0.0.1"},
		},
		{
			Name:  "multiple addresses",
			Input: []string{"127.0.0.1", "127.0.0.2"},
		},
	}

	for _, test := range cases {
		t.Run(test.Name, func(t *testing.T) {
			controller := NewTransportController(test.Input)
			defer controller.Close()
		})
	}
}

// Check that receive message will receive the correct message each time and block when there are no more messages
func TestReceiveMessage(t *testing.T) {
	controller := NewTransportController([]string{"127.0.0.2", "127.0.0.3", "127.0.0.4", "127.0.0.5", "127.0.0.6", "127.0.0.7"})
	defer controller.Close()

	type testCase struct {
		Name     string
		Payloads [][]byte
		Amounts  []uint
		From     []string
		Delay    []time.Duration
		Expected uint
	}

	cases := []testCase{
		{
			Name:     "single end",
			Payloads: [][]byte{{0xff}},
			Amounts:  []uint{20},
			From:     []string{"127.0.0.7:5999"},
			Delay:    []time.Duration{time.Microsecond * 16},
			Expected: 5100,
		},
		{
			Name:     "multiple ends",
			Payloads: [][]byte{{0xff}, {0xff, 0xff}, {0xff, 0xff, 0xff}, {0xff, 0xff, 0xff, 0xff}, {0xff, 0xff, 0xff, 0xff, 0xff}},
			Amounts:  []uint{20, 20, 20, 20, 20},
			From:     []string{"127.0.0.2:5999", "127.0.0.3:5999", "127.0.0.4:5999", "127.0.0.5:5999", "127.0.0.6:5999"},
			Delay:    []time.Duration{time.Microsecond * 8, time.Microsecond * 12, time.Microsecond * 16, time.Microsecond * 20, time.Microsecond * 24},
			Expected: 76500,
		},
	}

	for _, test := range cases {
		t.Run(test.Name, func(t *testing.T) {
			for i, src := range test.From {
				go sendToTCP(src, fmt.Sprintf("127.0.0.1:%d", serverPort), test.Payloads[i], test.Amounts[i], test.Delay[i], t)
			}

			payloadsChan := make(chan []byte, sumUint(test.Amounts))
			got := uint(0)
		loop:
			for {
				go func() {
					payloadsChan <- controller.ReceiveMessage()
				}()

				select {
				case payload := <-payloadsChan:
					for _, b := range payload {
						got += uint(b)
					}
				case <-time.After(time.Second):
					break loop
				}
			}

			if got != test.Expected {
				t.Errorf("expected sum of %v, got %v", test.Expected, got)
			}
		})
		controller.Close()
		serverPort++
		controller = NewTransportController([]string{"127.0.0.2", "127.0.0.3", "127.0.0.4", "127.0.0.5", "127.0.0.6", "127.0.0.7"})
	}
}

// Check that receive data will always return something and block when there is nothing to return
func TestReceiveData(t *testing.T) {

}

// Test that send will send exactly what it's told to send
func TestSend(t *testing.T) {

}

// Check that AliveConnections will only return the connections that are alive
func TestAliveConnections(t *testing.T) {

}

func dialTCP(src, dst string, t *testing.T) *net.TCPConn {
	srcAddr, err := net.ResolveTCPAddr("tcp", src)
	if err != nil {
		t.Fatal(err)
	}

	dstAddr, err := net.ResolveTCPAddr("tcp", dst)
	if err != nil {
		t.Fatal(err)
	}

	conn, err := net.DialTCP("tcp", srcAddr, dstAddr)
	if err != nil {
		fmt.Println(err)
		t.Fatal(err)
	}

	return conn
}

func sendToTCP(src, dst string, payload []byte, amount uint, delay time.Duration, t *testing.T) {
	<-time.After(delay)
	conn := dialTCP(src, dst, t)
	if conn == nil {
		return
	}
	defer conn.Close()

	for i := uint(0); i < amount; i++ {
		<-time.After(delay)

		n, err := conn.Write(payload)
		if err != nil {
			t.Error(err)
		} else {
			t.Log(n)
		}
	}
}

func sumUint(input []uint) (sum uint) {
	for _, n := range input {
		sum += n
	}
	return
}
