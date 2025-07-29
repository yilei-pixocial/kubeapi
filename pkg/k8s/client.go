package k8s

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"path/filepath"
)

func NewClientSet(kubeconfigPath string) (*kubernetes.Clientset, error) {
	// 首先尝试 InCluster 配置（如果在 Kubernetes 集群中运行）
	if config, err := rest.InClusterConfig(); err == nil {
		return kubernetes.NewForConfig(config)
	}

	// 否则使用 kubeconfig
	return NewClientSetFromConfig(kubeconfigPath)
}

func NewClientSetFromConfig(kubeconfigPath string) (*kubernetes.Clientset, error) {
	if kubeconfigPath == "" {
		kubeconfigPath = filepath.Join(homedir.HomeDir(), ".kube", "config")
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		return nil, err
	}

	return kubernetes.NewForConfig(config)
}
