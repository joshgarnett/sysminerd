package main

type CPUInputModule struct{}

func (m CPUInputModule) name() string {
	return "cpu"
}

func (m CPUInputModule) init(config map[interface{}]interface{}) error {
	return nil
}

func (m CPUInputModule) tearDown() error {
	return nil
}
