package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/x/errors"
	"github.com/sirupsen/logrus"
	"github.com/yilei-pixocial/kubeapi/pkg/k8s"
	"github.com/yilei-pixocial/kubeapi/pkg/sys/resp"
	"github.com/yilei-pixocial/kubeapi/pkg/sysinit"
	"github.com/yilei-pixocial/kubeapi/pkg/vo"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type K8sService struct {
	ClientSet kubernetes.Interface
}

func NewK8sService() (*K8sService, error) {

	clientSet, err := k8s.NewClientSetFromConfig(sysinit.GCF.UString("kubernetes.kubeconfig"))
	if err != nil {
		return nil, err
	}

	return &K8sService{
		ClientSet: clientSet,
	}, nil
}

// GetServices
// @Summary 获取集群service信息
// @Tags　k8s
// @Accept application/json
// @Produce application/json
// @Success 200 {object} model.Message
// @Failure 400 {object} model.Message
// @Router /api/v1/k8s/services [get]
func (k *K8sService) GetServices(ctx iris.Context) {

	namespaceList, err := k.ClientSet.CoreV1().Namespaces().List(ctx.Request().Context(), metav1.ListOptions{})
	if err != nil {
		ctx.JSON(resp.ErrorWithMsg(err.Error()))
		return
	}

	var results []vo.ServiceVo
	for _, ns := range namespaceList.Items {
		serviceList, err := k.ClientSet.CoreV1().Services(ns.Name).List(ctx.Request().Context(), metav1.ListOptions{})
		if err != nil {
			ctx.JSON(resp.ErrorWithMsg(err.Error()))
			return
		}

		for _, svc := range serviceList.Items {
			createAt := time.Unix(0, svc.CreationTimestamp.UnixMilli()*int64(time.Millisecond)).
				Format("2006-01-02 15:04:05")
			results = append(results, vo.ServiceVo{
				Name:       svc.Name,
				Namespace:  svc.Namespace,
				Type:       string(svc.Spec.Type),
				ClusterIP:  svc.Spec.ClusterIP,
				Ports:      svc.Spec.Ports,
				CreateTime: createAt,
				Selector:   svc.Spec.Selector,
			})
		}
	}

	err = ctx.JSON(resp.OkWithData(results))
	if err != nil {
		ctx.JSON(resp.ErrorWithMsg(err.Error()))
		return
	}
}

// GetNamespaces
// @Summary 获取集群命名空间信息
// @Tags　k8s
// @Accept application/json
// @Success 200 {object} model.Message
// @Failure 400 {object} model.Message
// @Router /api/v1/k8s/namespaces [get]
func (k *K8sService) GetNamespaces(ctx iris.Context) {

	namespaceList, err := k.ClientSet.CoreV1().Namespaces().List(ctx.Request().Context(), metav1.ListOptions{})
	if err != nil {
		ctx.JSON(resp.ErrorWithMsg(err.Error()))
		return
	}

	var results []vo.NamespaceVo
	for _, ns := range namespaceList.Items {
		if ns.Status.Phase == "Active" {
			results = append(results, vo.NamespaceVo{
				Name: ns.Name,
				CreateTime: time.Unix(0, ns.CreationTimestamp.UnixMilli()*int64(time.Millisecond)).
					Format("2006-01-02 15:04:05"),
			})
		}
	}

	err = ctx.JSON(resp.OkWithData(results))
	if err != nil {
		ctx.JSON(resp.ErrorWithMsg(err.Error()))
		return
	}
	return
}

func SyncToRedis() {

	kubeClient, err := k8s.NewClientSetFromConfig(sysinit.GCF.UString("kubernetes.kubeconfig"))
	if err != nil {
		logrus.Errorf("Failed to create kubernetes clientSet from config. Err: %+v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	setupSignalHandling(cancel)

	if err := syncData(ctx, kubeClient); err != nil {
		logrus.Errorf("Initial sync failed: %v", err)
	}

	ticker := time.NewTicker(60 * time.Second) // 每60秒同步一次
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := syncData(ctx, kubeClient); err != nil {
				logrus.Errorf("Sync failed: %v", err)
			}
		case <-ctx.Done():
			logrus.Errorf("Shutting down...")
			return
		}
	}
}

type ClusterData struct {
	ClusterId       string           `json:"clusterID"`
	ClusterName     string           `json:"clusterName"`
	ClusterRegionID string           `json:"clusterRegionID"`
	K8sVersion      string           `json:"k8sVersion"`
	Namespaces      []vo.NamespaceVo `json:"namespaces"`
	Services        []vo.ServiceVo   `json:"services"`
}

func syncData(ctx context.Context, kubeClient *kubernetes.Clientset) error {
	logrus.Info("Starting sync...")

	var clusterID, clusterName, k8sVersion, clusterRegionID string

	if os.Getenv("kubernetes.clusterID") != "" {
		clusterID = os.Getenv("kubernetes.clusterID")
	} else if sysinit.GCF.UString("kubernetes.clusterID") != "" {
		clusterID = sysinit.GCF.UString("kubernetes.clusterID")
	} else {
		return errors.New("kubernetes cluster ID is empty")
	}

	if os.Getenv("kubernetes.clusterName") != "" {
		clusterName = os.Getenv("kubernetes.clusterName")
	} else if sysinit.GCF.UString("kubernetes.clusterName") != "" {
		clusterName = sysinit.GCF.UString("kubernetes.clusterName")
	} else {
		return errors.New("kubernetes clusterName is empty")
	}

	if os.Getenv("kubernetes.clusterRegionID") != "" {
		clusterRegionID = os.Getenv("kubernetes.clusterRegionID")
	} else if sysinit.GCF.UString("kubernetes.clusterRegionID") != "" {
		clusterRegionID = sysinit.GCF.UString("kubernetes.clusterRegionID")
	} else {
		return errors.New("kubernetes clusterRegionID is empty")
	}

	version, err := kubeClient.Discovery().ServerVersion()
	if err != nil {
		logrus.Errorf("Failed to get Kubernetes server version: %v", err)
	} else {
		k8sVersion = version.GitVersion
	}

	namespaces, err := kubeClient.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		logrus.Errorf("Failed to list namespaces: %v", err)
		return fmt.Errorf("failed to list namespaces: %w", err)
	}

	var serviceVos []vo.ServiceVo
	for _, ns := range namespaces.Items {
		services, err := kubeClient.CoreV1().Services(ns.Name).List(ctx, metav1.ListOptions{})
		if err != nil {
			logrus.Errorf("Failed to list services in namespace %s: %v", ns.Name, err)
			continue
		}

		for _, svc := range services.Items {
			serviceVos = append(serviceVos, vo.ServiceVo{
				ServiceID: fmt.Sprintf("%s/%s/%s", clusterID, ns.Name, svc.Name),
				Name:      svc.Name,
				Namespace: svc.Namespace,
				Type:      string(svc.Spec.Type),
				ClusterIP: svc.Spec.ClusterIP,
				Ports:     svc.Spec.Ports,
				CreateTime: time.Unix(0, svc.CreationTimestamp.UnixMilli()*int64(time.Millisecond)).
					Format("2006-01-02 15:04:05"),
				Selector:    svc.Spec.Selector,
				ClusterName: clusterName,
				ClusterID:   clusterID,
				NamespaceID: fmt.Sprintf("%s/%s", clusterID, ns.Name),
			})
		}
	}

	var namespaceVos []vo.NamespaceVo
	for _, ns := range namespaces.Items {
		namespaceVos = append(namespaceVos, vo.NamespaceVo{
			NamespaceID: fmt.Sprintf("%s/%s", clusterID, ns.Name),
			Name:        ns.Name,
			CreateTime: time.Unix(0, ns.CreationTimestamp.UnixMilli()*int64(time.Millisecond)).
				Format("2006-01-02 15:04:05"),
			ClusterID:   clusterID,
			ClusterName: clusterName,
			Status:      string(ns.Status.Phase),
		})
	}

	data := ClusterData{
		ClusterId:       clusterID,
		ClusterName:     clusterName,
		ClusterRegionID: clusterRegionID,
		K8sVersion:      k8sVersion,
		Namespaces:      namespaceVos,
		Services:        serviceVos,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	key := sysinit.GCF.UString("redis.keyPrefix") + clusterID
	err = sysinit.RedisCli.Set(ctx, key, jsonData, 3600*time.Second).Err()
	if err != nil {
		return fmt.Errorf("failed to set data in Redis: %w", err)
	}

	logrus.Infof("Sync completed for cluster %s, stored in key %s", clusterName, key)
	return nil
}

func setupSignalHandling(cancel context.CancelFunc) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		log.Printf("Received signal: %v", sig)
		cancel()
	}()
}
