package dataservice

type Config struct {
	Queues struct {
		Push string `yaml:"push"`
	} `yaml:"queues"`
}
