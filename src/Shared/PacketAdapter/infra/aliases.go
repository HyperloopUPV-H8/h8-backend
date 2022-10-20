package infra

type IP string

func stringsToIPs(strings []string) (ips []IP) {
	ips = make([]IP, len(strings))
	for i, str := range strings {
		ips[i] = IP(str)
	}

	return
}

type Port uint16
