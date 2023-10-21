package observer

import "spysat/stream"

type Observer struct {
	Group   string                   `yaml:"group"`
	Streams map[string]stream.Stream `yaml:"streams"`
}
