package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

const devTokenHeader = "Development-Token"
const accessTokenHeader = "Authorization"

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
	URL        string
	StatusCode int
	Error      string
}

func main() {
	performPrechecks()

	accessLog := getAccessLog()
	replayAccessLogs(accessLog)
}

func replayAccessLogs(accessLog []Request) {
	var responseErrors []ResponseError = make([]ResponseError, 0)

	var httpClient = &http.Client{
		Timeout: time.Second * 2,
	}

	for _, request := range accessLog {
		url := "https://" + request.Source.APIBase + request.Source.Message
		if request.Source.Accept != "application/vnd.mendeley-community-document-list+json" {
			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				fmt.Println(err.Error())
			}
			req.Header.Add("Accept", request.Source.Accept)
			req.Header.Add(devTokenHeader, os.Getenv("devToken"))
			req.Header.Add(accessTokenHeader, os.Getenv("accessToken"))
			response, err := httpClient.Do(req)
			if err != nil {
				fmt.Printf("%v\t%v\n", request.Source.Message, err.Error())
			} else {
				fmt.Printf("%v\t%v\n", response.StatusCode, request.Source.Message)
			}
			time.Sleep(8 * time.Second)
		}
	}
	reportErrors(responseErrors)
}

func reportErrors(errors []ResponseError) {
	fmt.Println("--------------- Error log ---------------")
	for _, element := range errors {
		fmt.Printf("%v\t%v\t%v\n", element.URL, element.StatusCode, element.Error)
	}
}

func getAccessLog() []Request {
	raw, err := ioutil.ReadFile("./access-logs.json")
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

func performPrechecks() {
	if os.Getenv("accessToken") == "" || os.Getenv("devToken") == "" {
		fmt.Printf("Export accessToken and devToken into your env.\n")
		os.Exit(1)
	}
}
