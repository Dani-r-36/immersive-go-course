package buffer

import (
	"reflect"
	"testing"
)

func TestBytesReturn(t *testing.T){
	b := StringBuffer("hello world")
	got := b.Bytes()
	want := []byte("hello world")

	reflect.DeepEqual(got, want)

}

func TestBytesWrite(t *testing.T){
	b := StringBuffer("hello world")
	b.Write([]byte(" you cruel place"))
	got := b.Bytes()
	want := []byte("hello world you cruel place")

	reflect.DeepEqual(got, want)

}

func TestReadFit(t *testing.T){
	b := StringBuffer("hello world")
	sliceAmount := make([]byte, 11)
	remaining, read := b.Read(sliceAmount)
	reflect.DeepEqual(11,remaining)
	reflect.DeepEqual(b, read)

}

func TestReadSmall(t *testing.T){
	b := StringBuffer("hello world")
	sliceAmount := make([]byte, 5)
	remaining, read := b.Read(sliceAmount)
	reflect.DeepEqual(6,remaining)
	reflect.DeepEqual(read, ([]byte("hello")))
	
}