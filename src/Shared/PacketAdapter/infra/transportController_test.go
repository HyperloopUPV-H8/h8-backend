package infra

import (
	"fmt"
	"net"
	"reflect"
	"sort"
	"sync"
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
	}

	// Check that receive message will receive the total amount of bytes initially sent
	for _, test := range cases {
		t.Run(test.Name, func(t *testing.T) {
			controller := NewTransportController([]string{"127.0.0.2", "127.0.0.3", "127.0.0.4", "127.0.0.5", "127.0.0.6"})
			defer controller.Close()

			for i, src := range test.From {
				go sendTCP(src, fmt.Sprintf("127.0.0.1:%d", serverPort), test.Payloads[i], test.Amounts[i], test.Delay[i], t)
			}

			payloadsChan := make(chan [][]byte, sumUint(test.Amounts))
			got := uint(0)
		loop:
			for {
				go func() {
					payloadsChan <- controller.ReceiveMessages()
				}()

				select {
				case payloads := <-payloadsChan:
					for _, payload := range payloads {
						for _, b := range payload {
							got += uint(b)
						}
					}
				case <-time.After(time.Millisecond * 500):
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

		payloadsChan := make(chan [][]byte)
		go func() {
			payloadsChan <- controller.ReceiveMessages()
		}()

		select {
		case <-payloadsChan:
			t.Error("expected read to block")
		case <-time.After(time.Second):
		}
	})
	serverPort++

	// Check that a connection can connect midway through a read and still give results
	t.Run("hotplug", func(t *testing.T) {
		controller := NewTransportController([]string{"127.0.0.2", "127.0.0.3", "127.0.0.4", "127.0.0.5", "127.0.0.6"})
		defer controller.Close()

		payloadsChan := make(chan [][]byte)
		go func() {
			payloadsChan <- controller.ReceiveMessages()
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
			case <-time.After(time.Second):
				if !connected {
					sendTCP("127.0.0.2:5999", fmt.Sprintf("127.0.0.1:%d", serverPort), []byte{0xff}, 1, 0, t)
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
	type testCase struct {
		Name     string
		File     string
		Expected [][]byte
	}

	cases := []testCase{
		{
			Name:     "udp packets",
			File:     "tests/resources/udp.pcapng",
			Expected: [][]byte{{0xff}, {0xff}},
		},
		{
			Name:     "tcp packets",
			File:     "tests/resources/tcp.pcapng",
			Expected: [][]byte{},
		},
		{
			Name:     "mixed packets",
			File:     "tests/resources/mixed.pcapng",
			Expected: [][]byte{{0xff}, {0xff}, {0xff}},
		},
	}

	for _, test := range cases {
		t.Run(test.Name, func(t *testing.T) {
			snifferTarget = test.File
			snifferLive = false

			controller := NewTransportController([]string{"127.0.0.2", "127.0.0.3", "127.0.0.4", "127.0.0.5", "127.0.0.6"})
			defer controller.Close()

			payloadChan := make(chan []byte)
			i := 0
		loop:
			for {
				go func() {
					payloadChan <- controller.ReceiveData()
				}()

				select {
				case payload := <-payloadChan:
					if !reflect.DeepEqual(payload, test.Expected[i]) {
						t.Errorf("expected %v, got %v", payload, test.Expected[i])
					}
					i++
				case <-time.After(time.Second):
					if i != len(test.Expected) {
						t.Errorf("expected more data")
					}
					break loop
				}
			}
		})
		serverPort++
	}
}

func TestSend(t *testing.T) {
	type testCase struct {
		Name     string
		Payloads [][]byte
		Amounts  []uint
		From     []string
		Expected []uint
	}

	cases := []testCase{
		{
			Name:     "single end",
			Payloads: [][]byte{{0xff}},
			Amounts:  []uint{20},
			From:     []string{"127.0.0.2"},
			Expected: []uint{5100},
		},
		{
			Name:     "multiple ends",
			Payloads: [][]byte{{0xff}, {0xff, 0xff}, {0xff, 0xff, 0xff}, {0xff, 0xff, 0xff, 0xff}},
			Amounts:  []uint{10, 10, 10, 10, 10},
			From:     []string{"127.0.0.2", "127.0.0.3", "127.0.0.4", "127.0.0.5"},
			Expected: []uint{2550, 5100, 7650, 10200},
		},
	}

	// Check that receive message will receive the total amount of bytes initially sent
	for _, test := range cases {
		t.Run(test.Name, func(t *testing.T) {
			controller := NewTransportController([]string{"127.0.0.2", "127.0.0.3", "127.0.0.4", "127.0.0.5", "127.0.0.6"})
			defer controller.Close()

			var wait *sync.WaitGroup = &sync.WaitGroup{}
			for i, src := range test.From {
				wait.Add(1)
				go receiveTCP(fmt.Sprintf("%s:5999", src), fmt.Sprintf("127.0.0.1:%d", serverPort), len(test.Payloads[i]), test.Expected[i], time.Second, wait, t)
				<-time.After(time.Millisecond * 100)
				for j := uint(0); j < test.Amounts[i]; j++ {
					controller.Send(src, test.Payloads[i])
				}
			}

			wait.Wait()
		})
		serverPort++
	}

	t.Run("delete closed", func(t *testing.T) {
		controller := NewTransportController([]string{"127.0.0.2", "127.0.0.3", "127.0.0.4", "127.0.0.5", "127.0.0.6"})
		defer controller.Close()

		var wait *sync.WaitGroup = &sync.WaitGroup{}
		wait.Add(1)
		go receiveTCP("127.0.0.2:5999", fmt.Sprintf("127.0.0.1:%d", serverPort), 5, 1275, time.Second, wait, t)
		<-time.After(time.Millisecond * 100)
		controller.Send("127.0.0.2", []byte{0xff, 0xff, 0xff, 0xff, 0xff})
		if !reflect.DeepEqual(controller.AliveConnections(), []string{"127.0.0.2"}) {
			t.Error("connection should be alive")
		}
		wait.Wait()
		<-time.After(time.Second)

		controller.Send("127.0.0.2", []byte{0xff})
		<-time.After(time.Millisecond * 100)
		if !reflect.DeepEqual(controller.AliveConnections(), []string{}) {
			t.Error("connection should be dead")
		}
	})
	serverPort++
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
				sendTCP(src, fmt.Sprintf("127.0.0.1:%d", serverPort), []byte{}, 0, 0, t)
			}

			// Need to wait because the controller might or might not be busy when accepting the connection
			<-time.After(time.Millisecond * 700)
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
	conn.SetLinger(0)

	return conn
}

func sendTCP(src, dst string, payload []byte, amount uint, delay time.Duration, t *testing.T) {
	<-time.After(delay)
	conn := dialTCP(src, dst, t)
	if conn == nil {
		return
	}

	for i := uint(0); i < amount; i++ {
		<-time.After(delay)

		_, err := conn.Write(payload)
		if err != nil {
			t.Error(err)
		}
	}
}

func receiveTCP(src, dst string, packetLen int, expected uint, delay time.Duration, wait *sync.WaitGroup, t *testing.T) {
	defer wait.Done()

	conn := dialTCP(src, dst, t)
	if conn == nil {
		fmt.Println("no conn")
		return
	}
	defer func() {
		conn.Close()
	}()

	payloadsChan := make(chan []byte)

	got := uint(0)

loop:
	for {
		go func() {
			buf := make([]byte, packetLen)
			conn.Read(buf)
			payloadsChan <- buf
		}()

		select {
		case payload := <-payloadsChan:
			for _, b := range payload {
				got += uint(b)
			}
		case <-time.After(delay):
			break loop
		}
	}

	if got != expected {
		t.Errorf("expected %d, got %d", expected, got)
	}
}

func sumUint(input []uint) (sum uint) {
	for _, n := range input {
		sum += n
	}
	return
}
