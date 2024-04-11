package main

import (
	"bytes"
	"net/http"
	"testing"
	"time"
	"os"
	"log"
	"github.com/stretchr/testify/assert"
)

func TestTimeConversion(t *testing.T){
	t.Run("time now format test", func(t *testing.T){
		timeAfterDelay := time.Now().UTC().Add(time.Duration(5) * time.Second)
		httpTime := timeAfterDelay.Format(http.TimeFormat)
		// httpTime := "Thu, 11 Apr 2024 13:10:05 GMT"
		want := time.Duration(5) * time.Second
		got := timeConversion(httpTime)

		//Apparently go is too fast so can be off by nanoseconds (hour of debugging)
		tolerance := time.Duration(1000) * time.Millisecond
		var absValue time.Duration
		if (got-want) <0{ absValue = -(got-want)}else{absValue = got-want}
		if absValue>tolerance{
			t.Errorf("got %d want %d which had tolerance of %d given, %v", got, want, tolerance, httpTime)
		}
	})

	t.Run("time in normal format", func(t *testing.T){
		want := time.Duration(5) * time.Second
		got := timeConversion("5")
		wantRounded := int (want.Seconds() + 0.5)
		gotRounded := int (got.Seconds() + 0.5)
		//Apparently go is too fast so can be off by nanoseconds (hour of debugging)
		// fmt.Println(wantRounded, gotRounded)
		if gotRounded != wantRounded{
			t.Errorf("got %d want %d given, %v", gotRounded, wantRounded, 5)
		}
	})
}

func captureOutput(f func()) string {
    var buf bytes.Buffer
    log.SetOutput(&buf)
    f()
    log.SetOutput(os.Stderr)
    return buf.String()
}

func TestResponseCode(t *testing.T){
	response := &http.Response{
		Header:  http.Header{responseHeader: []string{"2"}},
	}

	output := captureOutput(func() {
		if !responseCode(response){
			t.Errorf("Expected responseCode to throw true and not False")
		}
	})
	assert.Equal(t, "Error: Sever busy, trying again in: 2s", output)
}