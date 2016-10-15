package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type Request struct {
	Source Source `json:"_source"`
	Id     string `json:"_id"`
}

type Source struct {
	Message string `json:"message"`
	Accept  string `json:"accept"`
	APIBase string `json:"host_header"`
}

type ResponseError struct {
	RequestLog Request
	StatusCode int
	Error      string
}

func main() {
	accessLog := getAccessLog()

	replayAccessLogs(accessLog)
}

func replayAccessLogs(accessLog []Request) {
	var responseErrors []ResponseError = make([]ResponseError, 0)

	var httpClient = &http.Client{
		Timeout: time.Second * 2,
	}

	for _, request := range accessLog {
		fmt.Printf("%s\n", request)
		url := "https://" + request.Source.APIBase + request.Source.Message
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			fmt.Println(err.Error())
		}
		req.Header.Add("Accept", request.Source.Accept)
		req.Header.Add("Development-Token", "NTQ1LDE0ODQyNTg2NjU5NzMselF0QmtPOXgyMDNzQTRDV0dKRm1wME0wU1Bn")
		req.Header.Add("Authorization", "Bearer MSwxNDc2NDg2MjY1NDI4LDQ2MzEyMjc4MSwyNzgsYWxsLCxodHRwczovL2ludGVybmFsLWRvY3MtbGl2ZS5tZW5kZWxleS5jb20sNGQ5My1lNmU4ZTI1Yzc1MTAzNjA2MmM1LTQxZmFmZWJjY2VjODM2ZTUsNDJhOGEyNTctYTkxMi0zYTFkLTg3N2MtYThjN2Q2ZDM5MTZhLDVKLXdXQ2pwM2ZCU1hMQkFuWnJWcEp4M0ZKRQ")
		response, err := httpClient.Do(req)
		if err != nil {
			fmt.Println(err.Error())
		}
		if !isSuccessCode(response.StatusCode) {
			responseErrors = append(responseErrors, ResponseError{request, response.StatusCode, "error"})
		}
		time.Sleep(1 * time.Second)
	}

	fmt.Println(responseErrors)
}

func getAccessLog() []Request {
	raw, err := ioutil.ReadFile("./test-access-log.json")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	accessLog := make([]Request, 0)
	json.Unmarshal(raw, &accessLog)
	return accessLog
}

func isSuccessCode(a int) bool {
	httpStatus := []int{200, 301, 302}
	for _, b := range httpStatus {
		if b == a {
			return true
		}
	}
	return false
}
