package myutil

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func Config(name string) map[string]string {
	fmt.Println("Reading Config: " + name)
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
