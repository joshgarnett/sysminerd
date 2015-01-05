package main

import (
	"errors"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

// Config stores all the config options for sysminerd
type Config struct {
	Interval   float64
	Hostname   string
	ConfigPath string `yaml:"config_path"`
}

type ModuleConfig struct {
	Name     string
	Enabled  bool
	Settings map[string]interface{}
}

func parseConfig(path string) Config {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("Error reading config: %v", err)
	}
	yamlConfig := Config{}

	err = yaml.Unmarshal(data, &yamlConfig)
	if err != nil {
		log.Fatalf("Error parsing yaml: %v", err)
	}

	log.Printf("Config: %v", yamlConfig)

	return yamlConfig
}

func parseModuleConfig(path string) ModuleConfig {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("Error reading config: %v", err)
	}
	moduleConfig := ModuleConfig{}
	err = yaml.Unmarshal(data, &moduleConfig)
	if err != nil {
		log.Fatalf("Error parsing yaml: %v", err)
	}

	return moduleConfig
}

func (config *ModuleConfig) SettingsString(key string) (string, error) {
	value, ok := config.Settings[key]
	if !ok {
		return "", errors.New("Key does not exist")
	}
	svalue, ok := value.(string)
	if !ok {
		return "", errors.New("value is not a string")
	}
	return svalue, nil
}

func (config *ModuleConfig) SettingsInt(key string) (int, error) {
	value, ok := config.Settings[key]
	if !ok {
		return 0, errors.New("Key does not exist")
	}
	ivalue, ok := value.(int)
	if !ok {
		return 0, errors.New("value is not an int")
	}
	return ivalue, nil
}

func (config *ModuleConfig) SettingsArray(key string) ([]interface{}, error) {
	value, ok := config.Settings[key]
	if !ok {
		return nil, errors.New("Key does not exist")
	}
	avalue, ok := value.([]interface{})
	if !ok {
		return nil, errors.New("value is not an array")
	}
	return avalue, nil
}

func (config *ModuleConfig) SettingsStringArray(key string) ([]string, error) {
	avalue, err := config.SettingsArray(key)
	if err != nil {
		return nil, err
	}
	values := make([]string, 0, len(avalue))
	for _, v := range avalue {
		sv, _ := v.(string)
		values = append(values, sv)
	}
	return values, nil
}
