package sniffer

import "testing"

var glob any

func BenchmarkSniffer(b *testing.B) {
	b.StopTimer()
	b.ResetTimer()
	sniffer := New("\\Device\\NPF_Loopback", true, DefaultConfig([]string{"127.0.0.2", "127.0.0.3"}, []string{"127.0.0.2", "127.0.0.3"}))
	defer sniffer.Close()

	var payload []byte
	sniffer.GetNext()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		payload, _ = sniffer.GetNext()
	}
	b.StopTimer()
	glob = payload
}
