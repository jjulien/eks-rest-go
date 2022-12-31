package rest

import (
	"context"
	"encoding/base64"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	"github.com/aws/aws-sdk-go-v2/service/eks/types"
	"github.com/jjulien/eks-rest-go/creds"
	"k8s.io/client-go/rest"
)

func DefaultConfig(ctx context.Context, clusterName string) (*rest.Config, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}
	return WithAwsConfig(ctx, clusterName, cfg)
}

func WithAwsConfig(ctx context.Context, clusterName string, cfg aws.Config) (*rest.Config, error) {
	clusterInfo, err := eksClusterInfo(ctx, clusterName, cfg)
	if err != nil {
		return nil, err
	}
	caPemData, err := base64.StdEncoding.DecodeString(*clusterInfo.CertificateAuthority.Data)
	if err != nil {
		return nil, err
	}
	tlsConfig := rest.TLSClientConfig{
		CAData: caPemData,
	}
	bearerToken, err := creds.BearerToken(ctx, clusterName, cfg)
	if err != nil {
		return nil, err
	}
	k8sConfig := &rest.Config{
		Host:            *clusterInfo.Endpoint,
		TLSClientConfig: tlsConfig,
		BearerToken:     bearerToken,
	}
	return k8sConfig, nil
}

func eksClusterInfo(ctx context.Context, clusterName string, cfg aws.Config) (*types.Cluster, error) {
	client := eks.NewFromConfig(cfg)
	input := eks.DescribeClusterInput{Name: &clusterName}
	response, err := client.DescribeCluster(ctx, &input)
	if err != nil {
		return nil, err
	}
	return response.Cluster, nil
}
