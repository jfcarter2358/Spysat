package window

type Window struct {
	Layout   string              `yaml:"layout"`
	Classes  map[string][]string `yaml:"classes"`
	Elements map[string][]string `yaml:"elements"`
}
