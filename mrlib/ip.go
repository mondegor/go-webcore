package mrlib

import (
	"encoding/binary"
	"errors"
	"net"
)

// IP2int - возвращает IP в виде числа или ошибку, если конвертация невозможна.
func IP2int(ip net.IP) (uint32, error) {
	if len(ip) == 16 {
		return 0, errors.New("no sane way to convert ipv6 into uint32")
	}

	return binary.BigEndian.Uint32(ip), nil
}

// MustIP2int - возвращает IP в виде числа или panic, если конвертация невозможна.
func MustIP2int(ip net.IP) uint32 {
	value, err := IP2int(ip)
	if err != nil {
		panic(err)
	}

	return value
}

// Int2ip - возвращает net.IP полученного из указанного целого числа.
func Int2ip(number uint32) net.IP {
	ip := make(net.IP, 4)

	binary.BigEndian.PutUint32(ip, number)

	return ip
}
