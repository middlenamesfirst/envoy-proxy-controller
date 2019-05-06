package main

import (
	"flag"
	"fmt"

	v1 "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
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

	watcher, err := clientset.CoreV1().Services(ns).Watch(listOptions)
	if err != nil {
		klog.Fatal(err)
	}

	ch := watcher.ResultChan()
	for event := range ch {
		service := event.Object.(*v1.Service)

		if service.Spec.Type != "NodePort" {
			continue
		}

		switch event.Type {
		case watch.Added:
			fmt.Printf("Added - %v:%v\n", service.Spec.ClusterIP, service.Spec.Ports[0].NodePort)
		case watch.Deleted:
			fmt.Printf("Deleted - %v:%v\n", service.Spec.ClusterIP, service.Spec.Ports[0].NodePort)
		}
	}
}
