package main

import (
	"testing"
	"reflect"
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
	remaining := b.Read(sliceAmount)
	reflect.DeepEqual(11,remaining)

}

func TestReadSmall(t *testing.T){
	b := StringBuffer("hello world")
	sliceAmount := make([]byte, 5)
	remaining := b.Read(sliceAmount)
	reflect.DeepEqual(6,remaining)
	
}