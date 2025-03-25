package util

import (
	"container/list"
	"encoding/json"
)

func Decode(src, target interface{}) error {
	jsonStr, err := json.Marshal(src)
	if err != nil {
		return err
	}
	return json.Unmarshal(jsonStr, target)
}

func List2Array(list *list.List) []interface{} {
	var len = list.Len()
	if len == 0 {
		return nil
	}
	var arr []interface{}
	for e := list.Front(); e != nil; e = e.Next() {
		arr = append(arr, e.Value)
	}
	return arr
}
