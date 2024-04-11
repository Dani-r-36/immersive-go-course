package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
	"strconv"
)

var url string = "http://localhost:8080/"
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

func responseCode(response *http.Response) bool {
	_, exists := response.Header[responseHeader]
	if exists {
		fmt.Println(response.Header[responseHeader][0])
		sleepTime := timeConversion(response.Header[responseHeader][0])
		time.Sleep(sleepTime)
		return true
	}
	return false
}

func main(){
	response, err := http.Get(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error making http request: %v\n", err)
		os.Exit(1)
	}
	defer response.Body.Close()// close response body 

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error in extracting data from body: %v\n", err)
		os.Exit(1)
	}

	text := string(body[:])

	fmt.Println(text)
	fmt.Printf("client: status code: %d\n", response.StatusCode)
	if response.StatusCode != 429 && response.StatusCode !=200{
		fmt.Fprintf(os.Stderr, "Error status code: %v\n", response.StatusCode)
		os.Exit(1)
	}
	if responseCode(response){
		main()
	}
	
}

