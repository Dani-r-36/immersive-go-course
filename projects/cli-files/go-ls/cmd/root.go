
package cmd

import (
	"fmt"
	"os/exec"
	"os"
    "io"
)

func ErrCommand(w io.Writer, err error){
    if err != nil {
        fmt.Fprintf(w,"Error here: %v", err)
        return
        }
}

func Execute() {
	cmd := "ls"
    out, err := exec.Command(cmd).Output()
    ErrCommand(os.Stderr, err)
    output := string(out[:])
    fmt.Println(output)
}