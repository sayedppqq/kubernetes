package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/retry"
	"os"
)

func main() {

	kubeconfig := flag.String("kubeconfig", "/home/appscodepc/.kube/config", "my kubeconfig file location")
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		// in cluster config
		config, err = rest.InClusterConfig()
		if err != nil {
			panic(fmt.Errorf("error while making config file"))
		}
	}

	clientset, err := dynamic.NewForConfig(config)
	if err != nil {
		panic(fmt.Errorf("error while making clientset"))
	}

	deploymentsClient := schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"}

	deployment := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "apps/v1",
			"kind":       "Deployment",
			"metadata": map[string]interface{}{
				"name": "demo-deployment",
			},
			"spec": map[string]interface{}{
				"replicas": 2,
				"selector": map[string]interface{}{
					"matchLabels": map[string]interface{}{
						"app": "demo",
					},
				},
				"template": map[string]interface{}{
					"metadata": map[string]interface{}{
						"labels": map[string]interface{}{
							"app":  "demo",
							"name": "demo-pod-deployment",
						},
					},

					"spec": map[string]interface{}{
						"containers": []map[string]interface{}{
							{
								"name":  "api-server",
								"image": "sayedppqq/api-server",
								"ports": []map[string]interface{}{
									{
										"name":          "http",
										"protocol":      "TCP",
										"containerPort": 80,
									},
								},
							},
						},
					},
				},
			},
		},
	}
	prompt()

	//deployment created
	fmt.Println("Creating deployment....")
	result, err := clientset.Resource(deploymentsClient).Namespace("default").Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		panic(fmt.Errorf("error while creating deployment"))
	}
	fmt.Printf("Created Deployment : %v \n", result.GetName())

	listing(clientset, deploymentsClient)
	prompt()

	//update deployment
	fmt.Println("Updating deployment....")
	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {

		result, err := clientset.Resource(deploymentsClient).Namespace("default").Get(context.TODO(), "demo-deployment", metav1.GetOptions{})
		if err != nil {
			panic(fmt.Errorf("failed to get latest version"))
		}

		unstructured.SetNestedField(result.Object, int64(1), "spec", "replicas")

		_, updateErr := clientset.Resource(deploymentsClient).Namespace("default").Update(context.TODO(), result, metav1.UpdateOptions{})
		return updateErr
	})

	if retryErr != nil {
		panic(fmt.Errorf("error while updating deployment"))
	}
	fmt.Println("deployment updated")
	listing(clientset, deploymentsClient)
	prompt()

	//deleting deployment
	deletePolicy := metav1.DeletePropagationForeground
	fmt.Println("Deleting deployment....")

	err = clientset.Resource(deploymentsClient).Namespace("default").Delete(context.TODO(), "demo-deployment", metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	})

	if err != nil {
		panic(fmt.Errorf("error while deleting"))
	}

	fmt.Println("deployment deleted")

	listing(clientset, deploymentsClient)
	prompt()

}
func int32Ptr(i int32) *int32 { return &i }

func listing(clientset *dynamic.DynamicClient, deploymentsClient schema.GroupVersionResource) {
	//listing deployment
	fmt.Println("\nListing deployment....")

	deployments, err := clientset.Resource(deploymentsClient).Namespace("default").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(fmt.Errorf("error while listing deployment"))
	}
	for _, deployment := range deployments.Items {
		fmt.Println(deployment.GetName())
	}

	//listing pods
	fmt.Println("\nListing pods....")
	pods, err := clientset.Resource(schema.GroupVersionResource{Group: "", Version: "v1", Resource: "pods"}).Namespace("default").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(fmt.Errorf("error while listing pods %s ", err))
	}
	for _, pod := range pods.Items {
		fmt.Println(pod.GetName())
	}
}
func prompt() {
	fmt.Printf("-> Press Enter to continue.")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		break
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	fmt.Println()
}
