package main

import (
	"fmt"
)

type Endpoint struct {
	ClusterIP string
	NodePort  uint32
}

func (e *Endpoint) String() string {
	return fmt.Sprintf("%s:%d", e.ClusterIP, e.NodePort)
}
