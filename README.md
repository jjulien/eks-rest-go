# EKS Rest Go
This module uses the [aws-sdk-go-v2](https://github.com/aws/aws-sdk-go-v2) module to lookup EKS cluster information,
and then uses the Kubernetes [client-go](https://github.com/kubernetes/client-go) module to create and return
a `*rest.Config` that is authenticated using IAM credentials.

This is useful when you need to connect to the Kubernetes Master API of an EKS cluster from outside the
cluster, such as from a Lambda or other AWS Service.

## Requirements
This module uses IAM credentials to describe the EKS Cluster, and requires the permission `eks:DescribeCluster` on
the resource you are trying to connect to.

#### Example IAM Policy
```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "AllowClusterLogin",
            "Effect": "Allow",
            "Action": "eks:DescribeCluster",
            "Resource": "arn:aws:eks:us-east-1:111122223333:cluster/my-cluster-name"
        }
    ]
}
```

This module also assumes you have configured an IAM User or Role with access to your cluster following the
AWS Guide [Enabling IAM user and role access to your cluster](https://docs.aws.amazon.com/eks/latest/userguide/add-user-role.html)

#### Example AWS Auth ConfigMag
```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  annotations:
  name: aws-auth
  namespace: kube-system
data:
  mapRoles: |
    - rolearn: arn:aws:iam::111122223333:/role/worker--node-role
      username: system:node:{{EC2PrivateDNSName}}
      groups:
        - system:bootstrappers
        - system:nodes
    - rolearn: arn:aws:iam::111122223333:role/my-custom-role
      username: arn:aws:iam::111122223333:role/my-custom-role
      groups:
        - my-custom-group
```

## Usage
Using a Default AWS Config (most common) - [Full Example](examples/defaultConfig/defaultConfig.go)
```go
import eksrest "github.com/jjulien/eks-rest-go/rest"

restConfig, _ := eksrest.DefaultConfig(context.TODO(), clusterName)
clientSet, _ := kubernetes.NewForConfig(restConfig)
```

Using a Custom AWS Config - [Full Example](examples/customConfig/customConfig.go)
```go
import eksrest "github.com/jjulien/eks-rest-go/rest"

defaultCfg, _ := config.LoadDefaultConfig(context.TODO())
stsClient := sts.NewFromConfig(defaultCfg)
credProvider := stscreds.NewAssumeRoleProvider(stsClient, customRoleArn)
awsCfg := aws.Config{
	Region:      "us-east-1",
	Credentials: credProvider,
}
restConfig, _ := eksrest.WithAwsConfig(context.TODO(). clusterName, awsCfg)
clientSet, _ := kubernetes.NewForConfig(restConfig)
```

## License
See [LICENSE](LICENSE)
