package vo

type ServiceVo struct {
	Name        string      `json:"name"`
	Namespace   string      `json:"namespace"`
	ClusterIP   string      `json:"clusterIP"`
	ClusterID   string      `json:"clusterID"`
	Ports       interface{} `json:"ports"`
	CreatedTime int64       `json:"createdTime"`
	ClusterName string      `json:"clusterName"`
	Selector    interface{} `json:"selector"`
	Type        string      `json:"type"`
}
