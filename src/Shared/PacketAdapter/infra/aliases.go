package PacketAdapter

type IP string

// stringsToIPs and ipsToStrings are really similar but due to go generics limitations
// they cannot be made into a single function because it doesnt support casting from one
// type to another and there isn't an interface that specifies such property
func stringsToIPs(strings []string) (ips []IP) {
	ips = make([]IP, len(strings))
	for i, str := range strings {
		ips[i] = IP(str)
	}

	return
}

func ipsToStrings(ips []IP) (strings []string) {
	strings = make([]string, len(ips))
	for i, ip := range ips {
		strings[i] = string(ip)
	}

	return
}

type Port uint16
