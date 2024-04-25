package cmd

import (
	"bytes"
	"fmt"
	"testing"
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

func TestCarryCommand(t *testing.T){
	t.Run("No filename passed", func(t *testing.T){
		buffer := bytes.Buffer{}
		carryCommand(&buffer, "ls", "")
		got := buffer.String()
		want := "root.go\nroot_test.go\n" 
		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
		
	})

	t.Run("go-ls ran in go-ls folder", func(t *testing.T){
		buffer := bytes.Buffer{}
		carryCommand(&buffer, "ls", "../")
		got := buffer.String()
		want := "assets\ncmd\ngo.mod\ngo.sum\nmain\nmain.go\n" 
		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
		
	})
}