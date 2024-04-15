package cmd

import (
	"testing"
	"fmt"
	"bytes"
)

func TestRootErr(t *testing.T){

		t.Run("Error exists", func(t *testing.T){
		buffer := bytes.Buffer{}
		errMsg := "test error"
		ErrCommand(&buffer, fmt.Errorf(errMsg))
		got := buffer.String()
		want := "Error here: test error" 
		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
		
	})

	t.Run("No error", func(t *testing.T){
		buffer := bytes.Buffer{}
		ErrCommand(&buffer, nil)
		got := buffer.String()
		want := "" 
		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
		
	})
}