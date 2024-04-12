
package cmd

import (
	"fmt"
	"os/exec"
	"os"
)



func Execute() {
	cmd := "ls"
    out, err := exec.Command(cmd).Output()

    if err != nil {
        fmt.Fprintf(os.Stderr,"Error: %v", err)
    }

    output := string(out[:])
    fmt.Println(output)
}