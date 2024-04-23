package main

import (
	"fmt"
	"net/http"
	// "os"
	"output-and-error-handling/mock"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)


func Test200(t *testing.T){
	transport := mock.NewMockRoundTripper(t)
	transport.StubResponse(http.StatusOK, "some weather", nil)
	defer transport.AssertGotRightCalls()

	f := &WeatherFetcher{client: http.Client{
		Transport: transport,
	}}

	err, weather := f.Fetch("http://doesnotexist.com/")

	require.NoError(t, err)
	require.Equal(t, "some weather", weather)
}

func TestFetch429(t *testing.T) {
	transport := mock.NewMockRoundTripper(t)
	headers := make(http.Header)
	headers.Set("retry-after", "1")
	transport.StubResponse(http.StatusTooManyRequests, "server overloaded", headers)
	defer transport.AssertGotRightCalls()

	f := &WeatherFetcher{client: http.Client{
		Transport: transport,
	}}

	start := time.Now()
	err, _ := f.Fetch("http://doesnotexist.com/")
	elapsed := time.Since(start)

	require.Equal(t, fmt.Errorf("busy, waited"), err)
	require.GreaterOrEqual(t, elapsed, 1*time.Second)
}

// func TestTimeConversion(t *testing.T){
// 	t.Run("time now format test", func(t *testing.T){
// 		timeAfterDelay := time.Now().UTC().Add(time.Duration(5) * time.Second)
// 		httpTime := timeAfterDelay.Format(http.TimeFormat)
// 		// httpTime := "Thu, 11 Apr 2024 13:10:05 GMT"
// 		want := time.Duration(5) * time.Second
// 		got := timeConversion(httpTime)

// 		//Apparently go is too fast so can be off by nanoseconds (hour of debugging)
// 		tolerance := time.Duration(1000) * time.Millisecond
// 		var absValue time.Duration
// 		if (got-want) <0{ absValue = -(got-want)}else{absValue = got-want}
// 		if absValue>tolerance{
// 			t.Errorf("got %d want %d which had tolerance of %d given, %v", got, want, tolerance, httpTime)
// 		}
// 	})

// 	t.Run("time in normal format", func(t *testing.T){
// 		want := time.Duration(5) * time.Second
// 		got := timeConversion("5")
// 		wantRounded := int (want.Seconds() + 0.5)
// 		gotRounded := int (got.Seconds() + 0.5)
// 		//Apparently go is too fast so can be off by nanoseconds (hour of debugging)
// 		// fmt.Println(wantRounded, gotRounded)
// 		if gotRounded != wantRounded{
// 			t.Errorf("got %d want %d given, %v", gotRounded, wantRounded, 5)
// 		}
// 	})
// }
