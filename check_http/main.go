package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

type moduleArgs struct {
	Timeout int     `json:"timeout"`
	Checks  []check `json:"checks"`
}

type check struct {
	Name             string            `json:"name"`
	Method           string            `json:"method"`
	URL              string            `json:"url"`
	NoFollowRedirect bool              `json:"no_follow_redirect"`
	Headers          map[string]string `json:"headers"`
	Expected         expected          `json:"expected"`
}

type expected struct {
	StatusCode int               `json:"status_code"`
	Headers    map[string]string `json:"headers"`
}

type response struct {
	Checks  map[string]string `json:"checks"`
	Changed bool              `json:"changed"`
	Failed  bool              `json:"failed"`
}

func (c check) checkHTTP(timeout int) (bool, string) {
	if c.Method == "" {
		c.Method = "GET"
	}
	req, err := http.NewRequest(c.Method, c.URL, nil)
	if err != nil {
		jsonExit(true, map[string]string{"msg": err.Error()})
	}

	for k, v := range c.Headers {
		req.Header.Add(k, v)
		if k == strings.ToLower("host") {
			req.Host = v
		}
	}

	rFunc := func(r *http.Request,
		via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	if !c.NoFollowRedirect {
		rFunc = nil
	}

	client := &http.Client{
		Timeout:       time.Second * time.Duration(timeout),
		CheckRedirect: rFunc,
	}

	resp, err := client.Do(req)
	if err != nil {
		jsonExit(true, map[string]string{"msg": err.Error()})
	}

	passed := true

	statusMsg := ""
	// check status code
	if resp.StatusCode == c.Expected.StatusCode {
		statusMsg = fmt.Sprintf("status code: %v matches %v\n", resp.StatusCode, c.Expected.StatusCode)
	} else {
		statusMsg = fmt.Sprintf("status code: %v does not match %v\n", resp.StatusCode, c.Expected.StatusCode)
		passed = false
	}
	// check headers
	headerMsg := ""
	if len(c.Expected.Headers) > 0 {
		headerMsg = "headers: values match"
		for k, v := range c.Expected.Headers {
			av := resp.Header.Get(k)
			headerMsg = fmt.Sprintf("headers: %s expected %s got %s", k, v, av)
			if av != v {
				passed = false
			}
		}
	}

	return passed, statusMsg + headerMsg
}

func jsonExit(failed bool, msg map[string]string) {
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

func main() {
	if len(os.Args) != 2 {
		jsonExit(true, map[string]string{"msg": "incorrect num of args"})
	}

	bs, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		jsonExit(true, map[string]string{"msg": err.Error()})
	}

	m := &moduleArgs{}

	err = json.Unmarshal(bs, m)
	if err != nil {
		jsonExit(true, map[string]string{"msg": err.Error()})
	}

	if len(m.Checks) < 1 {
		jsonExit(true, map[string]string{"msg": "checks is empty"})
	}

	if m.Timeout <= 0 {
		m.Timeout = 3
	}

	results := make(map[string]string)
	var failed bool

	var wg sync.WaitGroup
	for _, c := range m.Checks {
		wg.Add(1)
		go func(c check) {
			passed, msg := c.checkHTTP(m.Timeout)
			if !passed {
				failed = true
			}
			results[c.Name] = msg
			wg.Done()
		}(c)
	}

	wg.Wait()
	jsonExit(failed, results)

}
