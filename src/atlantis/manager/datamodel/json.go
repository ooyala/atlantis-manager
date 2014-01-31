/* Copyright 2014 Ooyala, Inc. All rights reserved.
 *
 * This file is licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
 * except in compliance with the License. You may obtain a copy of the License at
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License is
 * distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and limitations under the License.
 */

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
