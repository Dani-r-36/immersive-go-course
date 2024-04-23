package main


//create new type, interface or struct because were extending bytes methods
type BytesSlice struct {
	//slices are dynamic so size depends on computers memory
	data []byte
	pos int
}

//converts string into byte in our struct
func StringBuffer(word string) BytesSlice{
	return BytesSlice{[]byte(word), 0}
}

//all methods below can't be called in goroutines because causes synchronization

func (b *BytesSlice) Bytes()[]byte{
	return b.data
}


//its extension... as we are appending a slice and not an element
func (b *BytesSlice) Write(extension []byte){
	b.data = append(b.data, extension...)
}

func(b *BytesSlice) Read(index []byte)int{
	remaining := len(b.data) - b.pos
	total := len(index)
	if remaining < total{
		total = remaining
	}
	copy(index, b.data[b.pos:b.pos+total])
	b.pos += total
	return total
}