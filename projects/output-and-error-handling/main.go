package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

type WeatherFetcher struct {
	client http.Client
}

var brokenRetryWord string = "a while"
var whileTime int = 7
var responseHeader string = "Retry-After"



func timeConversion(waitTime string) time.Duration {
	passedTime, err := time.Parse(http.TimeFormat, waitTime)
	if err == nil{
		sleepTime := passedTime.Sub(time.Now().UTC())
		fmt.Fprintf(os.Stderr, "Error: server busy, trying again in: %v\n", sleepTime)
		return sleepTime
	}else{
		num, err := strconv.Atoi(waitTime)
		if err!=nil {
			fmt.Fprintf(os.Stderr, "Error: couldn't convert sleep time into int giving default time. Error received: %v\n", err)
			num = whileTime
		}
		sleepTime := time.Duration(num) * time.Second
		fmt.Fprintf(os.Stderr, "Error: Sever busy, trying again in: %v\n", sleepTime)
		return sleepTime
	}
}

func responseCode(response *http.Response) (error ,bool) {
	_, exists := response.Header[responseHeader]
	if exists {
		sleepTime := timeConversion(response.Header[responseHeader][0])
		time.Sleep(sleepTime)
		return nil , exists
	}
	return fmt.Errorf("error handling too many requests"), exists
}

func(w *WeatherFetcher) Fetch(url string) (error, string){
	//pass w the struct http client and not http, xo that we can mock the client for testing 
	response, err := w.client.Get(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error making http request: %v\n", err)
		return err, ""
	}
	defer response.Body.Close()// close response body 
	switch response.StatusCode {
	case http.StatusOK:
		body, err := io.ReadAll(response.Body)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error in extracting data from body: %v\n", err)
			return err, ""
		}
		return nil, string(body)
	case http.StatusTooManyRequests:
		fmt.Fprintf(os.Stderr, "Error status code: %v\n", response.StatusCode)
		err, waited:=responseCode(response)
		if waited != false{
			return fmt.Errorf("busy, waited"), ""
		}
		return err, ""
	default:
		return fmt.Errorf("Error status code: %v\n",response.StatusCode), ""
	}
}

func main(){
	f := WeatherFetcher{}
	for {
		err, data := f.Fetch("http://localhost:8080/")
		if err != nil && err != fmt.Errorf("busy, waited"){
			os.Exit(1)
		} else if data != "" {
			fmt.Fprintf(os.Stdout, "%s\n", data)
			break
		}
	}
}

