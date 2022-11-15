package sniffer

import (
	"reflect"
	"testing"
	"time"
)

func TestSniffer(t *testing.T) {
	type testCase struct {
		file     string
		expected [][]byte
	}

	tests := map[string]testCase{
		"udp": {
			file:     "tests/resources/udp.pcapng",
			expected: [][]byte{{0xff}, {0xfe}, {0xfd}, {0xfc}, {0xfb}, {0xfa}},
		},
		"tcp": {
			file:     "tests/resources/tcp.pcapng",
			expected: [][]byte{{0xff}, {0xfe}, {0xfd}, {0xfc}, {0xfb}, {0xfa}},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			sniffer := New(test.file, false, []Filterer{})

			for _, payload := range test.expected {
				if got := sniffer.GetNextValidPayload(); !reflect.DeepEqual(payload, got) {
					t.Fatalf("expected %v, got %v", payload, got)
				}
			}
		})
	}
}

func TestBlock(t *testing.T) {
	sniffer := New("tests/resources/empty.pcapng", false, []Filterer{})

	done := make(chan bool)
	go func() {
		sniffer.GetNextValidPayload()
		done <- true
	}()

	select {
	case <-time.After(time.Second):
		t.Log("sniffer blocks on read")
	case <-done:
		t.Fatalf("expected sniffer to block on read")
	}
}

func TestFilters(t *testing.T) {
	sniffer := New("tests/resources/mixed.pcapng", false, []Filterer{
		SourceIPFilter{[]string{"127.0.0.2", "127.0.0.3", "127.0.0.4"}},
		DestinationIPFilter{[]string{"127.0.0.1", "127.0.0.2", "127.0.0.3", "127.0.0.4"}},
		UDPFilter{},
	})

	expected := [][]byte{{0xff}, {0xfd}, {0xfc}}

	got := make([][]byte, 0, len(expected))
loop:
	for {
		payload := make(chan []byte)
		go func() {
			payload <- sniffer.GetNextValidPayload()
		}()

		select {
		case <-time.After(time.Millisecond * 10):
			break loop
		case data := <-payload:
			got = append(got, data)
		}
	}

	if !reflect.DeepEqual(expected, got) {
		t.Fatalf("expected %v, got %v", expected, got)
	}
}

func BenchmarkFilters(b *testing.B) {
	sniffer := New("tests/resources/mixed.pcapng", false, []Filterer{
		SourceIPFilter{[]string{"127.0.0.2", "127.0.0.3", "127.0.0.4"}},
		DestinationIPFilter{[]string{"127.0.0.1", "127.0.0.2", "127.0.0.3", "127.0.0.4"}},
		UDPFilter{}})

	packet, err := sniffer.source.NextPacket()
	if err != nil {
		b.Fatal(err)
	}

	for i := 0; i < b.N; i++ {
		sniffer.applyFilters(&packet)
	}
}
