package structs

// Environment is representation of a single environment configuration.
type Environment struct {
	Name             string
	Domain           string
	Foundations      []string `yaml:",flow"`
	Authenticate     bool
	SkipSSL          bool `yaml:"skip_ssl"`
	Instances        uint16
	DisableRollback  bool                   `yaml:"rollback_disabled"`
	CustomParams     map[string]interface{} `yaml:"custom_params"`
	AllowInvalidUser bool                   `yaml:"allow_invalid_user"`
}
