package aliases

type IP string

// stringsToIPs and ipsToStrings are really similar but due to go generics limitations
// they cannot be made into a single function because it doesnt support casting from one
// type to another and there isn't an interface that specifies such property
func StringsToIPs(strings []string) (ips []IP) {
	ips = make([]IP, len(strings))
	for i, str := range strings {
		ips[i] = IP(str)
	}

	return ips
}

func IPsToStrings(ips []IP) (strings []string) {
	strings = make([]string, len(ips))
	for i, ip := range ips {
		strings[i] = string(ip)
	}

	return strings
}

type Port uint16

type Payload []byte

func PayloadsToBytes(payloads []Payload) (bytes [][]byte) {
	bytes = make([][]byte, len(payloads))
	for i, payload := range payloads {
		bytes[i] = payload
	}

	return bytes
}
