package types

type DeploymentInfo struct {
	Namespace string `json:"namespace,omitempty"` //命名空间
	Name      string `json:"name,omitempty"`      //部署名称
	Status    string `json:"status,omitempty"`    //状态，0：正常，1：可用，2：错误
	ImageName string `json:"image_name"`          //镜像名称
	InsCount  string `json:"ins_count"`           //实例数，如‘1/3’，右边数字是定义的最大实例数，左边是实际在运行的实例	数
}

type NamespaceInfo struct {
	Total      int              `json:"total,omitempty"`
	DeployList []DeploymentInfo `json:"list,omitempty"`
}

type QueryDeployment struct {
	Namespace string `json:"namespace,omitempty"` //命名空间
	Name      string `json:"name,omitempty"`      //部署名称
	Start     int    `json:"start,omitempty"`     //分页索引
	Limit     int    `json:"limit,omitempty"`     //分页条数，0表示查询全部
}

type QueryDeploymentOut struct {
	Total      int             `json:"total,omitempty"`
	Namespaces []NamespaceInfo `json:"namespace_list,omitempty"`
}

type QueryOut struct {
	Total int              `json:"total,omitempty"`
	List  []DeploymentInfo `json:"list,omitempty"`
}
