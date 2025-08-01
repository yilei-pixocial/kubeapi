package vo

type ServiceVo struct {
	ServiceID   string      `json:"serviceID"`
	Name        string      `json:"name"`
	Namespace   string      `json:"namespace"`
	ClusterIP   string      `json:"clusterIP"`
	ClusterID   string      `json:"clusterID"`
	Ports       interface{} `json:"ports"`
	CreateTime  string      `json:"createdTime"`
	ClusterName string      `json:"clusterName"`
	Selector    interface{} `json:"selector"`
	Type        string      `json:"type"`
	NamespaceID string      `json:"namespaceID"`
}
