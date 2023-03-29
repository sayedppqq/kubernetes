package main

import (
	"context"
	"flag"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	kubeconfig := flag.String("kubeconfig", "/home/appscodepc/.kube/config", "my kubeconfig file location")

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)

	if err != nil {

		// in cluster config
		config, err = rest.InClusterConfig()
		if err != nil {
			fmt.Printf("error while building config file %s", err.Error())
		}
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Printf("error %s while building clientset", err.Error())
	}
	cnt := 0
	for {
		cnt++
		fmt.Println(cnt, ">>>>>>>>>>>>>>>>>>>>>>")
		pods, err := clientset.CoreV1().Pods("default").List(context.Background(), metav1.ListOptions{})
		if err != nil {
			fmt.Printf(err.Error())
		}
		for _, pod := range pods.Items {
			fmt.Println(pod.Name)
		}
		clientset.AppsV1().Deployments()
	}
}
