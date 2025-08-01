package vo

type NamespaceVo struct {
	NamespaceID string `json:"namespaceID"`
	Name        string `json:"name"`
	CreateTime  string `json:"createdTime"`
	ClusterID   string `json:"clusterID"`
	ClusterName string `json:"clusterName"`
	Status      string `json:"status"`
}
