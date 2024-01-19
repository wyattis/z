package znet

import (
	"fmt"
	"net"
	"time"
)

/*
 * Iterate over port range to find an open port using net.Dial. If timeout == 0, the default timeout will be used
 */
func FindPortDial(host string, min, max int, timeout time.Duration) (port int, err error) {
	if timeout == 0 {
		timeout = time.Millisecond * 50
	}
	for port = min; port <= max; port++ {
		addr := net.JoinHostPort(host, fmt.Sprint(port))
		conn, err := net.DialTimeout("tcp", addr, timeout)
		if err != nil {
			return port, nil
		}
		if conn != nil {
			conn.Close()
		}
	}
	err = fmt.Errorf("Failed to find open port in range %d-%d", min, max)
	return
}

/*
 * Find N unique open ports by calling net.Listen with ":0" to let the OS assign a port
 */
func FindNPortsListen(n int) (ports []int, err error) {
	for i := 0; i < n; i++ {
		l, err := net.Listen("tcp", ":0")
		if err != nil {
			return ports, err
		}
		ports = append(ports, l.Addr().(*net.TCPAddr).Port)
		defer l.Close()
	}
	return
}

/**
 * Find an open port by calling net.Listen with ":0" to let the OS assign a port
 */
func FindPortListen() (port int, err error) {
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}

// Get preferred outbound ip of this machine
func GetOutboundIP() (net.IP, error) {
	conn, err := net.Dial("udp", "8.7.6.5:69")
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP, nil
}
