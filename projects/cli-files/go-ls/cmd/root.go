
package cmd

import (
	"fmt"
	"os/exec"
	"os"
    "io"
)

func ErrCommand(w io.Writer, err error) error{
    if err != nil {
        fmt.Fprintf(w,"Error here: %v", err)
        return err
        }
    return nil
}

func carryCommand(w io.Writer,cmd string, filename string) error{
    var out []byte
    var err error
    if filename == ""{out, err = exec.Command(cmd).Output()}else{out, err = exec.Command(cmd, filename).Output()}
    err = ErrCommand(os.Stderr, err)
    if err != nil{return err}
    output := string(out[:])
    fmt.Fprintf(w, output)
    return nil
}

func Execute() {
	cmd := "ls"
    if len(os.Args)>1{
        for _, fileName:= range os.Args[1:]{
            err := carryCommand(os.Stderr,cmd,fileName)
            if err != nil {
                fmt.Fprintln(os.Stderr, err)
                os.Exit(1)
                }
        }
        return
    }
    err := carryCommand(os.Stderr,cmd,"")
        if err != nil {
            fmt.Fprintln(os.Stderr, err)
            os.Exit(1)
            }
}