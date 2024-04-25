package main

import (
	"bytes"
	"reflect"
	"testing"
)

func TestWriterReturn(t *testing.T){
	writerTests := []struct {
		got    []string
		want   string
	}{
		{got: []string{"Rectangle"}, want: "Rectangle"},
		{got: []string{"2Circle"}, want: "Circle"},
		{got: []string{"start=2 end=10"}, want: "start= end="},
	}
	for _, tt := range writerTests {
		t.Run(tt.want, func(t *testing.T) {
			//existing method to convert bytes to string
			b := bytes.NewBufferString("")
			//passed buffer into writer
			filteredpipe := NewFilteringPipe(b)
			for _, inputData := range tt.got{
				n, _ := filteredpipe.WriteNext([]byte(inputData))
				reflect.DeepEqual(len(inputData),n)
			}
		})

	}

}


