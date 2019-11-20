package structs

type ErrorMatcherDescriptor struct {
	Description string `yaml:"description"`
	Pattern     string `yaml:"pattern"`
	Solution    string `yaml:"solution"`
	Code        string `yaml:"code"`
}
