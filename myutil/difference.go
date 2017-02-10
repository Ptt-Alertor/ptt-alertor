package myutil

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

func DifferenceJSON(jsonOld string, jsonNew string) []byte {
	mapOld, mapNew := jsonToMap(jsonOld), jsonToMap(jsonNew)
	mapDiff := DifferenceMap(mapOld, mapNew)

	jsonDiff, _ := json.Marshal(mapDiff)

	return jsonDiff
}

func DifferenceMap(mapOld []map[string]string, mapNew []map[string]string) []map[string]string {
	mapDiff := make([]map[string]string, 0)
	for _, objectNew := range mapNew {
		for index, objectOld := range mapOld {
			if reflect.DeepEqual(objectNew, objectOld) {
				break
			}
			if index == len(mapNew)-1 {
				mapDiff = append(mapDiff, objectNew)
			}
		}
	}
	return mapDiff
}

func jsonToMap(jsonString string) []map[string]string {
	var articles []map[string]string
	err := json.NewDecoder(strings.NewReader(jsonString)).Decode(&articles)
	if err != nil {
		fmt.Println(err)
	}
	return articles
}
