package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	eksrest "github.com/jjulien/eks-rest-go/rest"
	"k8s.io/client-go/kubernetes"
)

const clusterName = "my-cluster-name"
const customRoleArn = "my-role-arn"

func main() {
	clientSet, err := CustomConfigExample()
	if err != nil {
		panic(err)
	}
	fmt.Println(clientSet.ServerVersion())
}

func CustomConfigExample() (*kubernetes.Clientset, error) {
	defaultCfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}
	stsClient := sts.NewFromConfig(defaultCfg)
	credProvider := stscreds.NewAssumeRoleProvider(stsClient, customRoleArn)
	awsCfg := aws.Config{
		Region:      "us-east-1",
		Credentials: credProvider,
	}
	restConfig, err := eksrest.WithAwsConfig(context.TODO(), clusterName, awsCfg)
	if err != nil {
		return nil, err
	}
	clientSet, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, err
	}
	return clientSet, nil
}
