package tcp

import (
	"log"
	"net"
	"sync"
	"testing"
	"time"
)

func TestTCP(t *testing.T) {
	t.Run("echo client", func(t *testing.T) {
		tcp := Open(&Config{
			LocalPort:   5000,
			RemoteIPs:   []string{"127.0.0.2"},
			RemotePorts: []uint16{5001},
			Snaplen:     32,
		})
		defer tcp.Close()

		wait := sync.WaitGroup{}
		wait2 := sync.WaitGroup{}
		wait3 := sync.WaitGroup{}
		wait4 := sync.WaitGroup{}
		wait2.Add(1)
		wait3.Add(1)
		wait.Add(1)
		wait4.Add(1)

		tcp.SetOnRead(func(b []byte) {
			log.Println("server", string(b))
			wait4.Done()
		})

		go func() {
			defer wait.Done()
			laddr, err := net.ResolveTCPAddr("tcp", "127.0.0.2:5001")
			if err != nil {
				log.Fatalln(err)
			}
			raddr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:5000")
			if err != nil {
				log.Fatalln(err)
			}
			conn, err := net.DialTCP("tcp", laddr, raddr)
			if err != nil {
				log.Fatalln(err)
			}
			defer conn.Close()
			wait3.Done()
			wait2.Wait()
			buf := make([]byte, 4)
			n, err := conn.Read(buf)
			if err != nil {
				log.Fatalln(n, err)
			}
			log.Println("client", string(buf))
			_, err = conn.Write(buf[:n])
			if err != nil {
				log.Fatalln(err)
			}
		}()

		wait3.Wait()
		<-time.After(time.Second)
		err := tcp.Send("127.0.0.2", []byte("hi!\n"))
		if err != nil {
			t.Fatal(err)
		}
		wait2.Done()
		wait.Wait()
		wait4.Wait()
	})

	t.Run("echo server", func(t *testing.T) {
		tcp := Open(&Config{
			LocalPort:   7000,
			RemoteIPs:   []string{"127.0.0.2"},
			RemotePorts: []uint16{7001},
			Snaplen:     32,
		})
		defer tcp.Close()

		wait := sync.WaitGroup{}
		wait2 := sync.WaitGroup{}
		wait3 := sync.WaitGroup{}
		wait4 := sync.WaitGroup{}
		wait2.Add(1)
		wait3.Add(1)
		wait.Add(1)
		wait4.Add(1)

		tcp.SetOnRead(func(b []byte) {
			log.Println("client", string(b))
			wait4.Done()
		})

		go func() {
			defer wait.Done()
			laddr, err := net.ResolveTCPAddr("tcp", "127.0.0.2:7001")
			if err != nil {
				log.Fatalln(err)
			}
			listener, err := net.ListenTCP("tcp", laddr)
			if err != nil {
				log.Fatalln(err)
			}
			conn, err := listener.Accept()
			if err != nil {
				log.Fatalln(err)
			}
			defer conn.Close()
			wait3.Done()
			wait2.Wait()
			buf := make([]byte, 4)
			n, err := conn.Read(buf)
			if err != nil {
				log.Fatalln(n, err)
			}
			log.Println("server", string(buf))
			_, err = conn.Write(buf[:n])
			if err != nil {
				log.Fatalln(err)
			}
		}()

		wait3.Wait()
		<-time.After(time.Second)
		err := tcp.Send("127.0.0.2", []byte("hi!\n"))
		if err != nil {
			t.Fatal(err)
		}
		wait2.Done()
		wait.Wait()
		wait4.Wait()
	})
}
