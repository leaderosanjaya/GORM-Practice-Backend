package remoteconfig

type DefaultValue struct {
	Value string `json:"value"`
}

type Parameter struct {
	DefaultValue DefaultValue `json:"defaultValue"`
}

type Config struct {
	Parameters map[string]Parameter `json:"parameters"`
}
