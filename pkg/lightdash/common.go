package lightdash

type DbtConnection struct {
	Type           string `json:"type"`
	Repository     string `json:"repository"`
	Branch         string `json:"branch"`
	ProjectSubPath string `json:"project_sub_path"`
	HostDomain     string `json:"host_domain"`
	PersonalAccessToken string `json:"personal_access_token,omitempty"`
}

type WarehouseConnection struct {
	Type                   string `json:"type"`
	Account                string `json:"account,omitempty"`
	Role                   string `json:"role,omitempty"`
	Database               string `json:"database,omitempty"`
	Warehouse              string `json:"warehouse,omitempty"`
	Schema                 string `json:"schema,omitempty"`
	ClientSessionKeepAlive bool   `json:"clientSessionKeepAlive,omitempty"`
	Threads                int    `json:"threads,omitempty"`
	ServerHostName         string `json:"serverHostName,omitempty"`
	HTTPPath               string `json:"httpPath,omitempty"`
	PersonalAccessToken    string `json:"personalAccessToken,omitempty"`
	Catalog                string `json:"catalog,omitempty"`
}
