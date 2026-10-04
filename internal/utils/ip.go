package utils

import (
	"encoding/binary"
	"errors"
	"net"
)

func IPv4ToInt(ipv4 string) (uint32, error) {
	ip := net.ParseIP(ipv4)
	if ip == nil {
		return 0, errors.New("invalid IPv4 address")
	}
	ipv4Bytes := ip.To4()
	if ipv4Bytes == nil {
		return 0, errors.New("not a valid IPv4 address")
	}
	return binary.BigEndian.Uint32(ipv4Bytes), nil
}

func IntToIPv4(ip uint32) net.IP {
	ipBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(ipBytes, ip)
	return net.IP(ipBytes)
}
