package main

import (
	"fmt"
	"sync"

	v2 "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	"github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	"github.com/envoyproxy/go-control-plane/envoy/api/v2/endpoint"
	"github.com/envoyproxy/go-control-plane/pkg/cache"
)

type Logger interface {
	Infof(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

type Datastore struct {
	Endpoints     map[string]*Endpoint
	SnapshotCache cache.SnapshotCache
	Logger        Logger
	RWMutex       *sync.RWMutex
}

func NewDatastore(snapshotCache cache.SnapshotCache, logger Logger) *Datastore {
	return &Datastore{
		Endpoints:     map[string]*Endpoint{},
		SnapshotCache: snapshotCache,
		Logger:        logger,
		RWMutex:       &sync.RWMutex{},
	}
}

func (d *Datastore) SetSnapshot() {
	endpoints := []cache.Resource{}
	clusters := []cache.Resource{}
	routes := []cache.Resource{}
	listeners := []cache.Resource{}

	d.RWMutex.RLock()
	for _, endpoint := range d.Endpoints {
		e := MakeEndpoint("", endpoint.ClusterIP, endpoint.NodePort)
		endpoints = append(endpoints, e)
	}
	d.RWMutex.RUnlock()

	snapshot := cache.NewSnapshot("1.0", endpoints, clusters, routes, listeners)
	_ = d.SnapshotCache.SetSnapshot("node1", snapshot)
}

func (d *Datastore) AddEndpoint(clusterIP string, nodePort int32) {
	// Defers are last in first out. We don't want to create the snapshot until the mutex is unlocked
	defer d.SetSnapshot()

	d.RWMutex.Lock()
	defer d.RWMutex.Unlock()

	addr := d.makeMapKey(clusterIP, nodePort)

	d.Endpoints[addr] = &Endpoint{
		ClusterIP: clusterIP,
		NodePort:  uint32(nodePort),
	}

	d.Logger.Infof("Added endpoint %s", addr)
}

func (d *Datastore) DeleteEndpoint(clusterIP string, nodePort int32) {
	// Defers are last in first out. We don't want to create the snapshot until the mutex is unlocked
	defer d.SetSnapshot()

	d.RWMutex.Lock()
	defer d.RWMutex.Unlock()

	addr := d.makeMapKey(clusterIP, nodePort)
	delete(d.Endpoints, addr)

	d.Logger.Infof("Deleted endpoint %s", addr)
}

func (d *Datastore) makeMapKey(clusterIP string, nodePort int32) string {
	return fmt.Sprintf("%s:%d", clusterIP, nodePort)
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

func MakeEndpoint(clusterName string, address string, port uint32) *v2.ClusterLoadAssignment {
	return &v2.ClusterLoadAssignment{
		ClusterName: clusterName,
		Endpoints: []endpoint.LocalityLbEndpoints{{
			LbEndpoints: []endpoint.LbEndpoint{{
				HostIdentifier: &endpoint.LbEndpoint_Endpoint{
					Endpoint: &endpoint.Endpoint{
						Address: &core.Address{
							Address: &core.Address_SocketAddress{
								SocketAddress: &core.SocketAddress{
									Protocol: core.TCP,
									Address:  address,
									PortSpecifier: &core.SocketAddress_PortValue{
										PortValue: port,
									},
								},
							},
						},
					},
				},
			}},
		}},
	}
}
