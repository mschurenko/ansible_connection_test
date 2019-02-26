package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strconv"
	"sync"
	"time"
)

type ModuleArgs struct {
	Timeout int     `json:"timeout"`
	Checks  []Check `json:"checks"`
}

type Check struct {
	Name     string   `json:"name"`
	Host     string   `json:"host"`
	Port     int      `json:"port"`
	Expected Expected `json:expected`
}

type Expected struct {
	Open bool `json:"open"`
}

type Response struct {
	Checks  map[string]string `json:"checks"`
	Changed bool              `json:"changed"`
	Failed  bool              `json:"failed"`
}

func (c Check) checkTCPPort(timeout int) bool {
	fmt.Println(c.Host + ":" + strconv.Itoa(c.Port))

	conn, err := net.DialTimeout("tcp", c.Host+":"+strconv.Itoa(c.Port), time.Second*time.Duration(timeout))
	if err != nil {
		return false
	}

	conn.Close()
	return true
}

func jsonExit(failed bool, msg map[string]string) {
	bs, _ := json.Marshal(Response{
		Checks:  msg,
		Changed: false,
		Failed:  failed,
	})
	fmt.Println(string(bs))
	os.Exit(0)
}

func main() {
	if len(os.Args) != 2 {
		jsonExit(true, map[string]string{"msg": "incorrect num of args"})
	}

	bs, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		jsonExit(true, map[string]string{"msg": err.Error()})
	}

	m := &ModuleArgs{}

	err = json.Unmarshal(bs, m)
	if err != nil {
		jsonExit(true, map[string]string{"msg": err.Error()})
	}

	if len(m.Checks) < 1 {
		jsonExit(true, map[string]string{"msg": "checks is empty"})
	}

	if m.Timeout <= 0 {
		m.Timeout = 1
	}

	results := make(map[string]string)
	var failed bool

	var wg sync.WaitGroup
	for _, check := range m.Checks {
		wg.Add(1)
		go func(check Check) {
			isOpen := check.checkTCPPort(m.Timeout)
			if isOpen && check.Expected.Open {
				results[check.Name] = "is OPEN and expected to be OPEN"
				failed = false
			} else if isOpen && !check.Expected.Open {
				results[check.Name] = "is OPEN and expected to be CLOSED"
				failed = true
			} else if !isOpen && !check.Expected.Open {
				results[check.Name] = "is CLOSED and expected to be CLOSED"
				failed = false
			} else if !isOpen && check.Expected.Open {
				results[check.Name] = "is CLOSED and expected to be OPEN"
				failed = true
			} else {
				results[check.Name] = "unknown response"
				failed = true
			}
			wg.Done()
		}(check)
	}

	wg.Wait()
	jsonExit(failed, results)

}
