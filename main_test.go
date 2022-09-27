package main

import (
	"fmt"
	"testing"
)

func TestBiliClient_GetHostList(t *testing.T) {
	hostList := TestBiliClient.GetHostList()
	fmt.Println(hostList)
	if len(hostList) != 4 {
		t.Errorf("Wrong number returned, expected 4, got %d", len(hostList))
	}
}
