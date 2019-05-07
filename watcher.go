package main

import (
	"context"
	"log"
	"sync"

	v1 "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
)

type Watcher struct {
	Namespace   string
	ListOptions meta_v1.ListOptions
	Clientset   *kubernetes.Clientset
	Datastore   *Datastore
}

func NewWatcher(ns string, options meta_v1.ListOptions, clientset *kubernetes.Clientset, ds *Datastore) *Watcher {
	return &Watcher{
		Namespace:   ns,
		ListOptions: options,
		Clientset:   clientset,
		Datastore:   ds,
	}
}

func (w *Watcher) Run(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)

	go func() {
		defer wg.Done()

		watcher, err := w.Clientset.CoreV1().Services(w.Namespace).Watch(w.ListOptions)
		if err != nil {
			log.Fatal(err)
		}

		ch := watcher.ResultChan()

		for {
			select {
			case <-ctx.Done():
				return
			case event := <-ch:

				service := event.Object.(*v1.Service)

				if service.Spec.Type != "NodePort" {
					continue
				}

				switch event.Type {
				case watch.Added:
					w.Datastore.AddEndpoint(service.Spec.ClusterIP, service.Spec.Ports[0].NodePort)
				case watch.Deleted:
					w.Datastore.DeleteEndpoint(service.Spec.ClusterIP, service.Spec.Ports[0].NodePort)
				}
			}
		}

	}()
}
