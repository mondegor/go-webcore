package mrlib

import (
	"encoding/binary"
	"fmt"
	"net"

	"github.com/mondegor/go-webcore/mrlog"
)

func IP2int(ip net.IP) (uint32, error) {
	if len(ip) == 16 {
		return 0, fmt.Errorf("no sane way to convert ipv6 into uint32")
	}

	return binary.BigEndian.Uint32(ip), nil
}

func IP2intMust(ip net.IP) uint32 {
	value, err := IP2int(ip)
	if err != nil {
		mrlog.Default().Error().Err(err).Send()
		return 0
	}

	return value
}

func Int2ip(number uint32) net.IP {
	ip := make(net.IP, 4)

	binary.BigEndian.PutUint32(ip, number)

	return ip
}
