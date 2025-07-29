package service

import (
	"github.com/kataras/iris/v12"
	"github.com/yilei-pixocial/kubeapi/pkg/k8s"
	"github.com/yilei-pixocial/kubeapi/pkg/sys/resp"
	"github.com/yilei-pixocial/kubeapi/pkg/sysinit"
	"github.com/yilei-pixocial/kubeapi/pkg/vo"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
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
			results = append(results, vo.ServiceVo{
				Name:        svc.Name,
				Namespace:   svc.Namespace,
				Type:        string(svc.Spec.Type),
				ClusterIP:   svc.Spec.ClusterIP,
				Ports:       svc.Spec.Ports,
				CreatedTime: svc.CreationTimestamp.UnixMilli(),
				Selector:    svc.Spec.Selector,
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
				Name:        ns.Name,
				CreatedTime: ns.CreationTimestamp.UnixMilli(),
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
