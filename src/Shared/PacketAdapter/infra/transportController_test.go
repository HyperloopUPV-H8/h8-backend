package PacketAdapter

import (
	"fmt"
	"net"
	"reflect"
	"sort"
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

func TestReceiveMessage(t *testing.T) {
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
			From:     []string{"127.0.0.2:5999"},
			Delay:    []time.Duration{time.Microsecond * 16},
			Expected: 5100,
		},
		{
			Name:     "multiple ends",
			Payloads: [][]byte{{0xff}, {0xff, 0xff}, {0xff, 0xff, 0xff}, {0xff, 0xff, 0xff, 0xff}},
			Amounts:  []uint{20, 20, 20, 20, 20},
			From:     []string{"127.0.0.2:5999", "127.0.0.3:5999", "127.0.0.4:5999", "127.0.0.5:5999"},
			Delay:    []time.Duration{time.Microsecond * 8, time.Microsecond * 12, time.Microsecond * 16, time.Microsecond * 20, time.Microsecond * 24},
			Expected: 51000,
		},
		{
			Name:     "zero ends",
			Payloads: [][]byte{},
			Amounts:  []uint{},
			From:     []string{},
			Delay:    []time.Duration{},
			Expected: 0,
		},
	}

	// Check that receive message will receive the total amount of bytes initially sent
	for _, test := range cases {
		t.Run(test.Name, func(t *testing.T) {
			controller := NewTransportController([]string{"127.0.0.2", "127.0.0.3", "127.0.0.4", "127.0.0.5", "127.0.0.6"})
			defer controller.Close()

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
		serverPort++
	}

	t.Run("block on read", func(t *testing.T) {
		controller := NewTransportController([]string{"127.0.0.2", "127.0.0.3", "127.0.0.4", "127.0.0.5", "127.0.0.6"})
		defer controller.Close()

		payloadsChan := make(chan []byte)
		go func() {
			payloadsChan <- controller.ReceiveMessage()
		}()

		select {
		case <-payloadsChan:
			t.Error("expected read to block")
		case <-time.After(time.Second * 5):
			t.Log("blocked")
		}
	})
	serverPort++

	// Check that a connection can connect midway through a read and still give results
	t.Run("hotplug", func(t *testing.T) {
		controller := NewTransportController([]string{"127.0.0.2", "127.0.0.3", "127.0.0.4", "127.0.0.5", "127.0.0.6"})
		defer controller.Close()

		payloadsChan := make(chan []byte)
		go func() {
			payloadsChan <- controller.ReceiveMessage()
		}()

		connected := false
	loop:
		for {
			select {
			case <-payloadsChan:
				if !connected {
					t.Error("expected to block, got message")
				} else {
					break loop
				}
			case <-time.After(time.Second * 5):
				if !connected {
					sendToTCP("127.0.0.2:5999", fmt.Sprintf("127.0.0.1:%d", serverPort), []byte{0xff}, 1, 0, t)
					connected = true
				} else {
					t.Error("expected to receive message, got block on read")
					break loop
				}
			}
		}
	})
	serverPort++
}

func TestReceiveData(t *testing.T) {

}

func TestSend(t *testing.T) {

}

func TestAliveConnections(t *testing.T) {
	type testCase struct {
		Name      string
		Connected []string
		Expected  []string
	}

	cases := []testCase{
		{
			Name:      "single end",
			Connected: []string{"127.0.0.2:5999"},
			Expected:  []string{"127.0.0.2"},
		},
		{
			Name:      "multiple ends",
			Connected: []string{"127.0.0.2:5999", "127.0.0.3:5999", "127.0.0.4:5999", "127.0.0.5:5999", "127.0.0.6:5999"},
			Expected:  []string{"127.0.0.2", "127.0.0.3", "127.0.0.4", "127.0.0.5", "127.0.0.6"},
		},
		{
			Name:      "invalid ends",
			Connected: []string{"127.0.0.2:5999", "127.0.0.10:5999", "127.0.0.14:5999", "127.0.0.5:5999", "127.0.0.7:5999"},
			Expected:  []string{"127.0.0.2", "127.0.0.5"},
		},
	}

	for _, test := range cases {
		t.Run(test.Name, func(t *testing.T) {
			controller := NewTransportController([]string{"127.0.0.2", "127.0.0.3", "127.0.0.4", "127.0.0.5", "127.0.0.6"})
			defer controller.Close()

			for _, src := range test.Connected {
				sendToTCP(src, fmt.Sprintf("127.0.0.1:%d", serverPort), []byte{}, 0, 0, t)
			}

			// Need to wait because the controller might or might not be busy when accepting the connection
			<-time.After(time.Millisecond)
			sort.Strings(test.Expected)
			got := controller.AliveConnections()
			sort.Strings(got)
			if !reflect.DeepEqual(test.Expected, got) {
				t.Errorf("expected %v, got %v", test.Expected, got)
			}
		})
		serverPort++
	}
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

		_, err := conn.Write(payload)
		if err != nil {
			t.Error(err)
		}
	}
}

func sumUint(input []uint) (sum uint) {
	for _, n := range input {
		sum += n
	}
	return
}
