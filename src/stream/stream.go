package stream

type Stream struct {
	Headers   []string               `yaml:"headers"`
	Types     []string               `yaml:"types"`
	Probe     string                 `yaml:"probe"`
	Analyst   string                 `yaml:"analyst"`
	Arguments map[string]interface{} `yaml:"arguments"`
	Status    string                 `yaml:"status"`
}
