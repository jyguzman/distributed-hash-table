package bson

import "fmt"

type Reader struct {
	pos  int
	data []byte
}

func NewReader(data []byte) *Reader {
	return &Reader{0, data}
}

func (r *Reader) ReadField() (string, error) {
	t := Type(r.data[r.pos])
	r.pos++
	for r.data[r.pos] != byte(0) {
		r.pos++
	}
	r.pos++
	fmt.Println(t)
	return "", nil
}
