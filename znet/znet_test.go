package znet

import "testing"

func TestFindNPorts(t *testing.T) {
	ports, err := FindNPortsListen(5)
	if err != nil {
		t.Error(err)
	} else if len(ports) != 5 {
		t.Errorf("expected %d ports", 5)
	}
}
