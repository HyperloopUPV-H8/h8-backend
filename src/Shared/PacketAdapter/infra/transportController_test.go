package PacketAdapter

import (
	"fmt"
	"net"
	"reflect"
	"testing"
	"time"
)

// Check that NewTransportController wont panic when being called
func TestNewTransportController(t *testing.T) {
	t.Run("no addr", func(t *testing.T) {
		NewTransportController([]string{})
	})
	serverPort++

	t.Run("one addr", func(t *testing.T) {
		NewTransportController([]string{"127.0.0.1"})
	})
	serverPort++

	t.Run("one invalid addr", func(t *testing.T) {
		NewTransportController([]string{"abcd"})
	})
	serverPort++

	t.Run("invalid addr", func(t *testing.T) {
		NewTransportController([]string{"127.0.0.1", "abcd"})
	})
	serverPort++

	t.Run("duped", func(t *testing.T) {
		NewTransportController([]string{"127.0.0.1", "127.0.0.1"})
	})
	serverPort++

	t.Run("multiple addr", func(t *testing.T) {
		NewTransportController([]string{"127.0.0.1", "127.0.0.2"})
	})
	serverPort++
}

// Check that receive message will receive the correct message each time and block when there are no more messages
func TestReceiveMessage(t *testing.T) {
	controller := NewTransportController([]string{"127.0.0.1"})

	go func() {
		laddr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:5999")
		if err != nil {
			panic(err)
		}
		raddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("127.0.0.1:%d", serverPort))
		if err != nil {
			panic(err)
		}
		conn, err := net.DialTCP("tcp", laddr, raddr)
		if err != nil {
			panic(err)
		}

		for i := 0; i < 5; i++ {
			conn.Write([]byte{1, 2, 3, 4, 5})
			// 160us is the theoretical amount of time the backend should take to read
			// all the messages the pod sends at 30mb/s
			<-time.After(time.Microsecond * 160)
		}
		conn.Close()
	}()

	t.Run("receive multiple messages", func(t *testing.T) {
		expected := make([][]byte, 1)
		expected[0] = make([]byte, packetMaxLength)
		expected[0][0] = 1
		expected[0][1] = 2
		expected[0][2] = 3
		expected[0][3] = 4
		expected[0][4] = 5

		for i := 0; i < 5; i++ {
			if got := controller.ReceiveMessage(); got == nil {
				t.Errorf("Expected []byte, got nil")
			} else if !reflect.DeepEqual(got, expected) {
				t.Errorf("Expected %v, got %v", expected, got)
			}
		}
	})

	t.Run("block with no messages", func(t *testing.T) {
		done := make(chan bool)
		go func() {
			controller.ReceiveMessage()

			done <- true
		}()

		timeout := time.NewTimer(time.Second * 5)

		select {
		case <-done:
			t.Errorf("Expected blocking, got value")
		case <-timeout.C:
			break
		}
	})

	serverPort++
}

// Check that receive data will always return something and block when there is nothing to return
func TestReceiveData(t *testing.T) {
	snifferTarget = "tests/udptraffic.pcapng"
	snifferLive = false

	controller := NewTransportController([]string{"127.0.0.2", "127.0.0.3", "127.0.0.4"})

	t.Run("receive multiple messages", func(t *testing.T) {
		for i := 0; i < 5; i++ {
			if got := controller.ReceiveData(); got == nil {
				t.Errorf("Expected []byte, got nil")
			}
		}
	})

	t.Run("block with no messages", func(t *testing.T) {
		done := make(chan bool)
		go func() {
			controller.ReceiveData()

			done <- true
		}()

		timeout := time.NewTimer(time.Second * 5)

		select {
		case <-done:
			t.Errorf("Expected blocking, got value")
		case <-timeout.C:
			break
		}
	})

	serverPort++
}

// Test that send will send exactly what it's told to send
func TestSend(t *testing.T) {
	controller := NewTransportController([]string{"127.0.0.2"})

	done := make(chan bool)
	go func() {
		laddr, err := net.ResolveTCPAddr("tcp", "127.0.0.2:5999")
		if err != nil {
			panic(err)
		}
		raddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("127.0.0.1:%d", serverPort))
		if err != nil {
			panic(err)
		}
		conn, err := net.DialTCP("tcp", laddr, raddr)
		if err != nil {
			panic(err)
		}
		done <- true

		for i := 0; i < 5; i++ {
			buf := make([]byte, 5, 5)
			n, _ := conn.Read(buf)
			if n != 5 {
				t.Error("expected 5 bytes, got", n)
			}

			if !reflect.DeepEqual(buf, []byte{1, 2, 3, 4, 5}) {
				t.Errorf("expected %v, got %v", []byte{1, 2, 3, 4, 5}, buf)
			}
		}
		conn.Close()
		done <- true
	}()

	<-done

	for i := 0; i < 5; i++ {
		<-time.After(time.Microsecond * 160)
		controller.Send("127.0.0.2", []byte{1, 2, 3, 4, 5})
	}

	<-done

	controller.Send("127.0.0.2", []byte{1, 2, 3, 4, 5})

	serverPort++
}

// Check that AliveConnections will only return the connections that are alive
func TestAliveConnections(t *testing.T) {
	controller := NewTransportController([]string{"127.0.0.2", "127.0.0.3"})

	t.Run("no connections alive", func(t *testing.T) {
		if !reflect.DeepEqual(controller.AliveConnections(), []string{}) {
			t.Errorf("expected no connection to be alive")
		}
	})

	t.Run("one connection alive", func(t *testing.T) {
		done := make(chan bool)
		go func() {
			laddr, err := net.ResolveTCPAddr("tcp", "127.0.0.3:5999")
			if err != nil {
				panic(err)
			}
			raddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("127.0.0.1:%d", serverPort))
			if err != nil {
				panic(err)
			}
			conn, err := net.DialTCP("tcp", laddr, raddr)
			if err != nil {
				panic(err)
			}
			done <- true
			conn.Close()
		}()

		<-done
		<-time.After(time.Microsecond)
		if !reflect.DeepEqual(controller.AliveConnections(), []string{"127.0.0.3"}) {
			t.Errorf("expected one connection to be alive")
		}
	})

	t.Run("all connections alive", func(t *testing.T) {
		done := make(chan bool)
		go func() {
			laddr, err := net.ResolveTCPAddr("tcp", "127.0.0.2:5999")
			if err != nil {
				panic(err)
			}
			raddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("127.0.0.1:%d", serverPort))
			if err != nil {
				panic(err)
			}
			conn, err := net.DialTCP("tcp", laddr, raddr)
			if err != nil {
				panic(err)
			}
			done <- true
			conn.Close()
		}()

		<-done
		<-time.After(time.Microsecond)
		if !reflect.DeepEqual(controller.AliveConnections(), []string{"127.0.0.3", "127.0.0.2"}) {
			t.Errorf("expected two connections to be alive")
		}
	})
}
