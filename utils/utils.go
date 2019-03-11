package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// JSONexit ...
func JSONexit(failed bool, msg map[string]string) {
	type response struct {
		Checks  map[string]string `json:"checks"`
		Changed bool              `json:"changed"`
		Failed  bool              `json:"failed"`
	}
	rc := 0
	if failed {
		rc = 1
	}
	bs, _ := json.Marshal(response{
		Checks:  msg,
		Changed: false,
		Failed:  failed,
	})
	fmt.Println(string(bs))
	os.Exit(rc)
}

// ChkArgs ...
func ChkArgs() {
	if len(os.Args) != 2 {
		JSONexit(true, map[string]string{"msg": "incorrect num of args"})
	}
}

// GetParams ...
func GetParams() []byte {
	xb, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		JSONexit(true, map[string]string{"msg": err.Error()})
	}

	return xb
}
