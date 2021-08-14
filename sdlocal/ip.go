package sdlocal

import (
	"net"

	"github.com/gaorx/stardust3/sderr"
)

type IpPred func(net.Interface, net.IP) bool

var (
	IpExtractor = func(_ net.Interface, addr net.Addr) net.IP {
		switch v := addr.(type) {
		case *net.IPNet:
			return v.IP
		case *net.IPAddr:
			return v.IP
		default:
			return nil
		}
	}
)

func AllIp(preds ...IpPred) ([]net.IP, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, sderr.WithStack(err)
	}
	ips := make([]net.IP, 0)
	for _, iface := range ifaces {
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			ip := IpExtractor(iface, addr)
			if len(ip) > 0 {
				ok := true
				for _, pred := range preds {
					if pred != nil && !pred(iface, ip) {
						ok = false
						break
					}
				}
				if ok {
					ips = append(ips, ip)
				}
			}
		}
	}
	return ips, nil
}

func Ip(preds ...IpPred) (net.IP, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, sderr.WithStack(err)
	}
	for _, iface := range ifaces {
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			ip := IpExtractor(iface, addr)
			if len(ip) > 0 {
				ok := true
				for _, pred := range preds {
					if pred != nil && !pred(iface, ip) {
						ok = false
						break
					}
				}
				if ok {
					return ip, nil
				}
			}
		}
	}
	return nil, sderr.New("not found ip")
}

func AllIpStr(preds ...IpPred) []string {
	ips, err := AllIp(preds...)
	if err != nil {
		return nil
	}
	r := make([]string, 0, len(ips))
	for _, ip := range ips {
		r = append(r, ip.String())
	}
	return r
}

func IpStr(preds ...IpPred) string {
	ip, err := Ip(preds...)
	if err != nil {
		return ""
	}
	return ip.String()
}

// Ip4

func makeIp4Preds(preds []IpPred) []IpPred {
	r := make([]IpPred, 0, len(preds)+1)
	r = append(r, OnlyIp4)
	r = append(r, preds...)
	return r
}

func AllIp4(preds ...IpPred) ([]net.IP, error) {
	return AllIp(makeIp4Preds(preds)...)
}

func Ip4(preds ...IpPred) (net.IP, error) {
	return Ip(makeIp4Preds(preds)...)
}

func AllIp4Str(preds ...IpPred) []string {
	return AllIpStr(makeIp4Preds(preds)...)
}

func Ip4Str(preds ...IpPred) string {
	return IpStr(makeIp4Preds(preds)...)
}

func InnerIp4() string {
	return IpStr(OnlyIp4, NotLoopback)
}

// IpPred

func OnlyIp4(_ net.Interface, ip net.IP) bool {
	ip4 := ip.To4()
	return len(ip4) > 0
}

func IfaceNameIs(ifaceName string) func(net.Interface, net.IP) bool {
	return func(iface net.Interface, ip net.IP) bool {
		return iface.Name == ifaceName
	}
}

func IfaceNameIn(ifaceNames ...string) func(net.Interface, net.IP) bool {
	return func(iface net.Interface, ip net.IP) bool {
		for _, ifaceName := range ifaceNames {
			if iface.Name == ifaceName {
				return true
			}
		}
		return false
	}
}

func OnlyLoopback(_ net.Interface, ip net.IP) bool {
	return ip.IsLoopback()
}

var (
	autoOsPreds = map[string]func(net.Interface, net.IP) bool{
		// TODO: 支持更多的系统
		"linux":  IfaceNameIn("eth0", "eth1"),
		"darwin": IfaceNameIn("en5", "en4", "en3", "en2", "en1", "en0"),
	}
)

func NotLoopback(iface net.Interface, ip net.IP) bool {
	pred, ok := autoOsPreds[OS()]
	if !ok {
		panic(sderr.New("unknown os for find ip"))
	}
	return pred(iface, ip)
}
