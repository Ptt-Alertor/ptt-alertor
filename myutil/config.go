package myutil

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

func Config(name string) map[string]string {
	projectRoot := ProjectRootPath()
	configJSON, err := ioutil.ReadFile(projectRoot + "/config/" + name + ".json")
	if err != nil {
		log.Fatal(err)
	}

	var config map[string]string
	err = json.Unmarshal(configJSON, &config)
	if err != nil {
		log.Fatal(err)
	}

	return config
}
