package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
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

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(fmt.Errorf("error while making clientset"))
	}

	deploymentsClient := clientset.AppsV1().Deployments("default")

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{Name: "demo-deployment"},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(2),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "demo",
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name: "client-go-deployment-pod",
					Labels: map[string]string{
						"app": "demo",
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  "api-server",
							Image: "sayedppqq/api-server",
							Ports: []apiv1.ContainerPort{
								{
									Name:          "http",
									ContainerPort: 80,
									Protocol:      apiv1.ProtocolTCP,
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
	result, err := deploymentsClient.Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		panic(fmt.Errorf("error while creating deployment"))
	}
	fmt.Printf("Created Deployment : %v \n", result.GetObjectMeta().GetName())

	listing(clientset)
	prompt()

	//update deployment
	fmt.Println("Updating deployment....")
	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {

		result, err := deploymentsClient.Get(context.TODO(), "demo-deployment", metav1.GetOptions{})
		if err != nil {
			panic(fmt.Errorf("failed to get latest version"))
		}
		result.Spec.Replicas = int32Ptr(1)

		_, updateErr := deploymentsClient.Update(context.TODO(), result, metav1.UpdateOptions{})
		return updateErr
	})

	if retryErr != nil {
		panic(fmt.Errorf("error while updating deployment"))
	}
	fmt.Println("deployment updated")
	listing(clientset)
	prompt()

	//deleting deployment
	deletePolicy := metav1.DeletePropagationForeground
	fmt.Println("Deleting deployment....")

	err = deploymentsClient.Delete(context.TODO(), "demo-deployment", metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	})

	if err != nil {
		panic(fmt.Errorf("error while deleting"))
	}

	fmt.Println("deployment deleted")

	listing(clientset)
	prompt()

}
func int32Ptr(i int32) *int32 { return &i }

func listing(clientset *kubernetes.Clientset) {
	//listing deployment
	fmt.Println("\nListing deployment....")

	deployments, err := clientset.AppsV1().Deployments("default").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(fmt.Errorf("error while listing deployment"))
	}
	for _, deployment := range deployments.Items {
		fmt.Println(deployment.Name)
	}

	//listing pods
	fmt.Println("\nListing pods....")
	pods, err := clientset.CoreV1().Pods("default").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(fmt.Errorf("error while listing pods"))
	}
	for _, pod := range pods.Items {
		fmt.Println(pod.Name)
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
