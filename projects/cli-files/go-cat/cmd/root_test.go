package cmd 

import (
	"testing"
	// "io/ioutil"
	"bytes"
	"os"
	"io"
)

func TestRootEmpty(t *testing.T){
	t.Run("No file given", func(t *testing.T){
		var buf bytes.Buffer
		// Redirect stdout to the buffer
		oldStdout := os.Stdout
		// os.Stdout = &buf
		os.Stdout = io.Writer(&buf)
		// os.Stdout = stdout
		defer func() {
			// Reset stdout
			os.Stdout = oldStdout
		}()
		ExecuteCmd()
		output := buf.String()
		want := "Use args to cat file: go-cat <file-1>%" 
		if output != want {
			t.Errorf("expected \"%s\" got \"%s\"", want, output)
	}
	})
}