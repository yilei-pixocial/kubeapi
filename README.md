## kubeapi

kubeapi is a Python library that provides a simple and efficient way to interact with Kubernetes clusters using the Kubernetes API. It allows you to perform various operations such as creating, updating, deleting, and retrieving resources in a Kubernetes cluster.

### Features

- Easy-to-use interface for interacting with Kubernetes API
- Supports Kubernetes Namespace、Service resources

### Starting with kubeapi

```shell
swag init

go run main.go
```

### Running with Docker

```shell
docker run --network=host -v $HOME/.kube/config:/root/.kube/config:ro -v $(pwd)/configs:/root/configs docker.io/library/kubeapi:v0.0.1
```

