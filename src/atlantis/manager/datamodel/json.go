package datamodel

import (
	"encoding/json"
	"log"
)

func getJson(nodePath string, data interface{}) error {
	raw_data, _, err := Zk.Get(nodePath)
	if err != nil {
		log.Printf("Error getting data from node %s. Error: %s.", nodePath, err.Error())
		return err
	}
	if len(raw_data) == 0 {
		raw_data = "{}"
	}
	err = json.Unmarshal([]byte(raw_data), data)
	if err != nil {
		log.Printf("Error decoding json: %s", err.Error())
		log.Println(raw_data)
	}
	return err
}

func setJson(nodePath string, data interface{}) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error encoding json: %s", err.Error())
		return err
	}
	_, err = Zk.TouchAndSet(nodePath, string(bytes))
	return err
}
