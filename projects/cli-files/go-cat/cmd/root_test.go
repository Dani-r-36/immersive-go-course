package cmd

import (
	"testing"
	"fmt"
	"bytes"
	"os"
)

func TestRoot(t *testing.T){
	t.Run("No file given", func(t *testing.T){
		buffer := bytes.Buffer{}
		argsBackup := os.Args
        os.Args = []string{"cmd"}
        defer func() { os.Args = argsBackup }()
		NoArgs(&buffer)
		got := buffer.String()
		want := "Use args to cat file: go-cat <file-1>" 
		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
		
	})

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