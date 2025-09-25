package mrtype

import (
	"encoding/binary"
	"errors"
	"net"
)

type (
	// DetailedIP - подробный IP.
	DetailedIP struct {
		Real  net.IP
		Proxy net.IP
	}
)

// ToUint - comment method.
func (ip *DetailedIP) ToUint() (realIP, proxyIP uint32, err error) {
	realIP, err = IP2uint(ip.Real)
	if err != nil {
		return 0, 0, err
	}

	proxyIP, err = IP2uint(ip.Proxy)
	if err != nil {
		return 0, 0, err
	}

	return realIP, proxyIP, nil
}

// String - comment method.
func (ip *DetailedIP) String() string {
	if len(ip.Proxy) == 0 || ip.Proxy.IsUnspecified() {
		return ip.Real.String()
	}

	return ip.Real.String() + ", " + ip.Proxy.String()
}

// LoadUint - comment method.
func (ip *DetailedIP) LoadUint(realIP, proxyIP uint32) {
	if realIP > 0 {
		ip.Real = Uint2ip(realIP)
	} else {
		ip.Real = nil // net.IP{}
	}

	if proxyIP > 0 {
		ip.Proxy = Uint2ip(proxyIP)
	} else {
		ip.Proxy = nil // net.IP{}
	}
}

// IP2uint - возвращает IP в виде числа или ошибку, если конвертация невозможна.
func IP2uint(ip net.IP) (uint32, error) {
	if len(ip) == 0 {
		return 0, nil
	}

	if ip4 := ip.To4(); ip4 != nil {
		return binary.BigEndian.Uint32(ip4), nil
	}

	if len(ip) == 16 {
		return 0, errors.New("no sane way to convert ipv6 into uint32")
	}

	return 0, errors.New("ip is incorrect")
}

// Uint2ip - возвращает net.IP полученного из указанного целого числа.
func Uint2ip(number uint32) net.IP {
	ip := make(net.IP, 4)

	binary.BigEndian.PutUint32(ip, number)

	return ip
}
