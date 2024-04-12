package cmd

import (
	"fmt"
	"io"
	"os"
	"os/exec"
)


func readLines(file io.Reader){
    buffer := make([]byte, 1024)
    for {
        numOfBytes, err := file.Read(buffer)
        if numOfBytes >0 {
            os.Stdout.Write(buffer[:numOfBytes])
        }
        if err != nil{
            break
        }
    }
}

func cat(fileName string){
    //open file and check if any error
    file, err := os.Open(fileName)
    if err != nil {
        fmt.Fprintf(os.Stderr,"Error opening file: %v", err)
    }
    //closes files after function finished
    defer file.Close()
    readLines(file)
}


func ExecuteCmd() {
	cmd := "cat"
    if len(os.Args)< 2 {
        fmt.Fprintf(os.Stderr, "Use args to cat file: go-cat <file-1>")
        return
    }
    // for _, fileName:= range os.Args[1:]{
    //     cat(fileName)
    // }
    for _, fileName:= range os.Args[1:]{
        out, err := exec.Command(cmd,fileName).Output()
        if err != nil {
        fmt.Fprintf(os.Stderr,"Error here: %v", err)
        }

        output := string(out[:])
        fmt.Fprintf(os.Stdout, output)
    }
}