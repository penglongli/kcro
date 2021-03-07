package utils

import "net"

func IFaces() ([]string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	var result []string
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			return nil, err
		}
		if addrs == nil || len(addrs) == 0 {
			continue
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.String() == "127.0.0.1" {
				continue
			}
			if ip.To4() != nil {
				result = append(result, ip.String())
			}
		}
	}
	return result, nil
}
