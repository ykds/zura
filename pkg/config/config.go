package config

import (
	"encoding/json"
	"gopkg.in/yaml.v2"
	"os"
	"strings"
)

func InitConfig(path string, cfg interface{}) {
	b, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	if strings.HasSuffix(path, ".json") {
		err = json.Unmarshal(b, cfg)
		if err != nil {
			panic(err)
		}
	}
	if strings.HasSuffix(path, ".yaml") {
		err = yaml.Unmarshal(b, cfg)
		if err != nil {
			panic(err)
		}
	}
}
