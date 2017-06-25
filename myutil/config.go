package myutil

import (
	"encoding/json"
	"io/ioutil"
)

func Config(name string) map[string]string {
	projectRoot := ProjectRootPath()
	filePath := projectRoot + "/config/" + name + ".json"
	configJSON, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	var config map[string]string
	err = json.Unmarshal(configJSON, &config)
	if err != nil {
		panic(err)
	}

	return config
}
