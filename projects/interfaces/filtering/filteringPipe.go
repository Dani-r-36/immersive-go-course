package main

import (
	"io"
)

type filteringPipe struct {
	//slices are dynamic so size depends on computers memory
	writer	io.Writer
}

//converts string into byte in our struct
func NewFilteringPipe(w io.Writer) filteringPipe {
	return filteringPipe{w}
}
// func NewFilteringPipe(writer io.Writer) filteringPipe {
// 	return filteringPipe{
// 		writer: writer,
// 	}
// }


func (fp *filteringPipe) WriteNext(data []byte) (int, error){
	for i:= range data{
		if data[i] <='0'|| data[i]>='9'{
			_, err := fp.writer.Write(data[i:i+1])
			if err != nil{
				return i, err
			}
		}
	}
	return len(data), nil
}