package znet

import (
	"fmt"
	"net"
	"time"
)

func FindOpenPortWithTimeout(timeout time.Duration, host string, min, max int) (port int, err error) {
	for port = min; port <= max; port++ {
		addr := net.JoinHostPort(host, fmt.Sprint(port))
		conn, err := net.DialTimeout("tcp", addr, timeout)
		if err != nil {
			fmt.Printf("Open port at: %s\n", addr)
			fmt.Println(err)
			return port, nil
		}
		if conn != nil {
			conn.Close()
		}
	}
	err = fmt.Errorf("Failed to find open port in range %d-%d", min, max)
	return
}

func FindOpenPort(host string, min, max int) (port int, err error) {
	return FindOpenPortWithTimeout(time.Second, host, min, max)
}
