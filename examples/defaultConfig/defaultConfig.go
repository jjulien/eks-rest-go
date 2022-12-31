package main

import (
	"context"
	"fmt"
	eksrest "github.com/jjulien/eks-rest-go/rest"
	"k8s.io/client-go/kubernetes"
)

const clusterName = "my-cluster-name"

func main() {
	clientSet, err := DefaultExample()
	if err != nil {
		panic(err)
	}
	fmt.Println(clientSet.ServerVersion())
}

func DefaultExample() (*kubernetes.Clientset, error) {
	restConfig, err := eksrest.DefaultConfig(context.TODO(), clusterName)
	if err != nil {
		return nil, err
	}
	clientSet, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, err
	}
	return clientSet, nil
}
