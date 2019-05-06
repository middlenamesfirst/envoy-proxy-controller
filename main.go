package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"time"

	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog"
)

func main() {
	var kubeconfig, master string
	flag.StringVar(&kubeconfig, "kubeconfig", "", "absolute path to the kubeconfig file")
	flag.StringVar(&master, "master", "", "master url")
	flag.Parse()

	// creates the connection
	config, err := clientcmd.BuildConfigFromFlags(master, kubeconfig)
	if err != nil {
		klog.Fatal(err)
	}

	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		klog.Fatal(err)
	}

	// watch services
	var ns, label, field string
	flag.StringVar(&ns, "namespace", "", "namespace")
	flag.StringVar(&label, "label", "", "Label selector")
	flag.StringVar(&field, "field", "", "Field selector")

	listOptions := meta_v1.ListOptions{
		LabelSelector: label,
		FieldSelector: field,
	}

	ds := NewDatastore()

	watcher := NewWatcher(ns, listOptions, clientset, ds)

	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}
	watcher.Run(ctx, wg)

	go func(d *Datastore, c context.Context) {
		ticker := time.NewTicker(1 * time.Second)

		for {
			select {
			case <-c.Done():
				return
			case <-ticker.C:
				fmt.Println("=========")
				fmt.Println("Endpoints")
				fmt.Println("=========")

				for _, endpoint := range d.GetEndpoints() {
					fmt.Println(endpoint)
				}
			}
		}
	}(ds, ctx)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	<-stop

	cancel()

	wg.Wait()
}
