package myutil

import (
	"bytes"
	"encoding/json"
	"reflect"
)

func DiffJSON(old, new []byte) (diff []byte) {
	mapDiff := DiffMap(jsonToMap(old), jsonToMap(new))
	diff, _ = json.Marshal(mapDiff)
	return diff
}

func DiffMap(old, new []map[string]string) (diff []map[string]string) {
	for _, nv := range new {
		for k, ov := range old {
			if reflect.DeepEqual(nv, ov) {
				break
			}
			if k == len(old)-1 {
				diff = append(diff, nv)
			}
		}
	}
	return diff
}

func jsonToMap(data []byte) (m []map[string]string) {
	if err := json.NewDecoder(bytes.NewReader(data)).Decode(&m); err != nil {
		panic(err)
	}
	return m
}
