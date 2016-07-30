package util

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

func WriteJson(p string, d interface{}) {
	b, err := json.MarshalIndent(d, "", "  ")
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(p, b, os.ModePerm)
	if err != nil {
		panic(err)
	}
}
