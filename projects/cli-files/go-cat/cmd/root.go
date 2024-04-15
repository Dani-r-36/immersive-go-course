package cmd

import (
	"fmt"
	"io"
	"os"
	"os/exec"
)



func NoArgs(w io.Writer){
    if len(os.Args)< 2 {
        fmt.Fprintf(w, "Use args to cat file: go-cat <file-1>")
        return
    }
}

func ErrCommand(w io.Writer, err error){
    if err != nil {
        fmt.Fprintf(w,"Error here: %v", err)
        return
        }
}

func ExecuteCmd() {
	cmd := "cat"
    NoArgs(os.Stderr)
    for _, fileName:= range os.Args[1:]{
        out, err := exec.Command(cmd,fileName).Output()
        ErrCommand(os.Stderr, err)
        output := string(out[:])
        fmt.Fprintf(os.Stdout, output)
    }
}