package src

import (
	"bytes"
	"encoding/json"
	"fmt"
)

func JSONStringToObject(jsonStr string, obj interface{}) bool {
	err := json.Unmarshal([]byte(jsonStr), &obj)
	if err != nil {
		return false
	} else {
		return true
	}
}

func ObjectToJSONStringFormatted(conf interface{}) string {
	b, err := json.Marshal(conf)
	if err != nil {
		return fmt.Sprintf("%+v", conf)
	}
	var out bytes.Buffer
	err = json.Indent(&out, b, "", "    ")
	if err != nil {
		return fmt.Sprintf("%+v", conf)
	}
	return out.String()
}

func ObjectToJSONString(obj interface{}) string {
	b, err := json.Marshal(obj)
	if err != nil {
		return "{}"
	} else {
		return string(b)
	}
}

