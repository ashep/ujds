package authservice

type Permission struct {
	Read  bool `json:"read"`
	Write bool `json:"write"`
}

type Permissions map[string]Permission

type Entity struct {
	Id          string
	Permissions Permissions
	Note        string
}
