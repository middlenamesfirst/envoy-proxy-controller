package main

import (
	"fmt"
	"sync"
)

type Datastore struct {
	Endpoints map[string]*Endpoint
	RWMutex   *sync.RWMutex
}

func NewDatastore() *Datastore {
	return &Datastore{
		Endpoints: map[string]*Endpoint{},
		RWMutex:   &sync.RWMutex{},
	}
}

func (d *Datastore) GetEndpoints() []*Endpoint {
	d.RWMutex.RLock()
	defer d.RWMutex.RUnlock()

	endpoints := []*Endpoint{}

	for _, e := range d.Endpoints {
		endpoints = append(endpoints, e)
	}

	return endpoints
}

func (d *Datastore) AddEndpoint(clusterIP string, nodePort int32) {
	d.RWMutex.Lock()
	defer d.RWMutex.Unlock()

	addr := d.makeMapKey(clusterIP, nodePort)

	d.Endpoints[addr] = &Endpoint{
		ClusterIP: clusterIP,
		NodePort:  nodePort,
	}
}

func (d *Datastore) DeleteEndpoint(clusterIP string, nodePort int32) {
	d.RWMutex.Lock()
	defer d.RWMutex.Unlock()

	addr := d.makeMapKey(clusterIP, nodePort)
	delete(d.Endpoints, addr)
}

func (d *Datastore) makeMapKey(clusterIP string, nodePort int32) string {
	return fmt.Sprintf("%s:%d", clusterIP, nodePort)
}
