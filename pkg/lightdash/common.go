package lightdash

type DbtConnection struct {
	Type           string `json:"type"`
	Repository     string `json:"repository"`
	Branch         string `json:"branch"`
	ProjectSubPath string `json:"project_sub_path"`
	HostDomain     string `json:"host_domain"`
}

type WarehouseConnection struct {
	Type                   string `json:"type"`
	Account                string `json:"account"`
	Role                   string `json:"role"`
	Database               string `json:"database"`
	Warehouse              string `json:"warehouse"`
	Schema                 string `json:"schema"`
	ClientSessionKeepAlive bool   `json:"clientSessionKeepAlive"`
	Threads                int    `json:"threads"`
}
