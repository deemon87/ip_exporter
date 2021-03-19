package info

import (
	"log"
	"net"
)

type Iface struct {
	Name       string
	Mac        net.HardwareAddr
	Ipv4       []string
	Ipv6       []string
	IsPrivate  bool
	IsLoopback bool
}

func GetInterfacesInfo(netInterface []net.Interface) []Iface {
	listInterfaces := []Iface{}
	for _, netInt := range netInterface {
		n, err := netInt.Addrs()
		if err != nil {
			log.Fatal(err)
		}
		if n != nil {
			ipFour, isPrivateFour, isLoopbackFour := getIpv4(n)
			ipSix, isPrivateSix, isLoopbackSix := getIpv6(n)
			isPrivate := isPrivateFour || isPrivateSix
			isLoopback := isLoopbackFour || isLoopbackSix

			listInterfaces = append(listInterfaces, Iface{
				Name:       netInt.Name,
				Mac:        netInt.HardwareAddr,
				Ipv4:       ipFour,
				Ipv6:       ipSix,
				IsPrivate:  isPrivate,
				IsLoopback: isLoopback,
			})
		}
	}
	return listInterfaces
}

func getIpv4(netAddress []net.Addr) ([]string, bool, bool) {
	var isPrivate bool
	var isLoopback bool

	listIpv4 := []string{}
	for _, nAdd := range netAddress {
		res, _, err := net.ParseCIDR(nAdd.String())
		if err != nil {
			log.Fatal(err)
		}
		ip := net.ParseIP(res.String())
		if ip.To4() != nil {
			isLoopback = ip.IsLoopback()
			isPrivate = checkPrivate(ip)
			listIpv4 = append(listIpv4, ip.String())
		}
	}
	return listIpv4, isPrivate, isLoopback
}

func getIpv6(netAddress []net.Addr) ([]string, bool, bool) {
	var isPrivate bool
	var isLoopback bool
	var isLinkLocalUnicast bool

	listIpv6 := []string{}
	for _, nAdd := range netAddress {
		res, _, err := net.ParseCIDR(nAdd.String())
		if err != nil {
			log.Fatal(err)
		}
		ip := net.ParseIP(res.String())
		if ip.To4() == nil && ip.To16() != nil {
			isLoopback = ip.IsLoopback()
			isPrivate = checkPrivate(ip)
			isLinkLocalUnicast = ip.IsLinkLocalUnicast()
			ipv6 := ip.String()
			if isLinkLocalUnicast {
				ipv6 = ""
			}

			listIpv6 = append(listIpv6, ipv6)
		}
	}
	return listIpv6, isPrivate, isLoopback
}

func checkPrivate(ip net.IP) bool {
	if ip4 := ip.To4(); ip4 != nil {
		return ip4[0] == 10 ||
			(ip4[0] == 172 && ip4[1]&0xf0 == 16) ||
			(ip4[0] == 192 && ip4[1] == 168)
	}
	return len(ip) == net.IPv6len && ip[0]&0xfe == 0xfc
}
