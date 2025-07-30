package vo

type NamespaceVo struct {
	Name        string `json:"name"`
	CreatedTime int64  `json:"createdTime"`
	ClusterID   string `json:"clusterID"`
	ClusterName string `json:"clusterName"`
	Status      string `json:"status"`
}
