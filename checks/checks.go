package checks

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"../utils"
)

// Checker ...
type Checker interface {
	Run() bool
	GetName() string
	GetMsg() string
}

type check struct {
	Timeout int    `json:"timeout"`
	Name    string `json:"name"`
	Msg     string `json:"msg"`
}

// GetName ...
func (c check) GetName() string {
	return c.Name
}

// GetMsg ...
func (c check) GetMsg() string {
	return c.Msg
}

type httpExpected struct {
	StatusCode int               `json:"status_code"`
	Headers    map[string]string `json:"headers"`
}

// HTTPcheck ...
type HTTPcheck struct {
	check
	URL              string            `json:"url"`
	Headers          map[string]string `json:"headers"`
	NoFollowRedirect bool              `json:"no_follow_redirect"`
	Expected         httpExpected      `json:"expected"`
}

// Run ...
func (h *HTTPcheck) Run() bool {
	req, err := http.NewRequest(http.MethodGet, h.URL, nil)
	if err != nil {
		utils.JSONexit(true, map[string]string{"msg": err.Error()})
	}

	for k, v := range h.Headers {
		req.Header.Add(k, v)
		// host header is special
		if "host" == strings.ToLower(k) {
			req.Host = v
		}
	}

	rFunc := func(r *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	if !h.NoFollowRedirect {
		rFunc = nil
	}

	client := &http.Client{
		Timeout:       time.Second * time.Duration(h.Timeout),
		CheckRedirect: rFunc,
	}

	resp, err := client.Do(req)
	if err != nil {
		utils.JSONexit(true, map[string]string{"msg": err.Error()})
	}

	passed := true
	var statusMsg string
	var headerMsg string

	// check status code
	if resp.StatusCode == h.Expected.StatusCode {
		statusMsg = fmt.Sprintf("status code: %v matches %v\n", resp.StatusCode, h.Expected.StatusCode)
	} else {
		statusMsg = fmt.Sprintf("status code: %v does not match %v\n", resp.StatusCode, h.Expected.StatusCode)
		passed = false
	}

	// check headers
	if len(h.Expected.Headers) > 0 {
		headerMsg = "headers: values match"
		for k, v := range h.Expected.Headers {
			av := resp.Header.Get(k)
			headerMsg = fmt.Sprintf("headers: %s expected %s got %s", k, v, av)
			if av != v {
				passed = false
			}
		}
	}

	h.Msg = statusMsg + headerMsg
	return passed
}

// HTTP ...
type HTTP struct {
	Checks []*HTTPcheck `json:"checks"`
}

// NewHTTP ...
func NewHTTP() *HTTP {
	utils.ChkArgs()
	xb := utils.GetParams()
	h := &HTTP{}
	err := json.Unmarshal(xb, h)
	if err != nil {
		utils.JSONexit(true, map[string]string{"msg": err.Error()})
	}

	return h
}

type portExpected struct {
	Open bool `json:"open"`
}

// PortCheck ...
type PortCheck struct {
	check
	Host     string       `json:"host"`
	Port     int          `json:"port"`
	Expected portExpected `json:"expected"`
}

// Port ...
type Port struct {
	Checks []*PortCheck `json:"checks"`
}

// NewPort ...
func NewPort() *Port {
	utils.ChkArgs()
	xb := utils.GetParams()
	p := &Port{}
	err := json.Unmarshal(xb, p)
	if err != nil {
		utils.JSONexit(true, map[string]string{"msg": err.Error()})
	}

	return p
}

// Run ...
func (p *PortCheck) Run() bool {
	conn, err := net.DialTimeout("tcp", p.Host+":"+strconv.Itoa(p.Port), time.Second*time.Duration(p.Timeout))
	if conn == nil {
		utils.JSONexit(true, map[string]string{"msg": err.Error()})
	}

	conn.Close()

	var passed bool

	if err == nil && p.Expected.Open {
		p.Msg = "is OPEN and expected to be OPEN"
		passed = true
	} else if err == nil && !p.Expected.Open {
		p.Msg = "is OPEN and expected to be CLOSED"
		passed = false
	} else if err != nil && !p.Expected.Open {
		p.Msg = "is CLOSED and expected to be CLOSED"
		passed = true
	} else if err != nil && p.Expected.Open {
		p.Msg = "is CLOSED and expected to be OPEN"
		passed = false
	} else {
		p.Msg = "Unknown response"
		passed = false
	}

	return passed
}

// RunChecks ...
func RunChecks(chks []Checker) {
	outputs := make(chan Checker, len(chks))

	for _, chk := range chks {
		go func(chk Checker) {
			passed := chk.Run()
			results := map[string]string{
				chk.GetName(): chk.GetMsg(),
			}
			if !passed {
				utils.JSONexit(true, results)
			}

			outputs <- chk
		}(chk)
	}

	allMsgs := make(map[string]string)
	var ci Checker

	for i := 0; i < len(chks); i++ {
		ci = <-outputs
		allMsgs[ci.GetName()] = ci.GetMsg()
	}

	utils.JSONexit(false, allMsgs)
}
