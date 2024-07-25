package bson

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strconv"
)

type Reader struct {
	pos  int
	data []byte
}

func NewReader(data []byte) *Reader {
	return &Reader{0, data}
}

func (r *Reader) ReadDocument() (*RawD, error) {
	raw := &RawD{Pairs: make(map[string]*Raw)}
	size, err := r.ReadSize()
	if err != nil {
		return nil, err
	}
	raw.Size = size
	for r.pos < int(size)-1 {
		field, err := r.ReadField()
		if err != nil {
			return nil, err
		}
		rawVal, err := r.ReadValue(field.Type)
		if err != nil {
			return nil, err
		}
		raw.Pairs[field.Name] = rawVal
	}
	//fmt.Println("pairs:", raw.Pairs)

	return raw, nil
}

func (r *Reader) ReadArray() (*RawArray, error) {
	raw := make(RawArray, 0)
	size, err := r.ReadSize()
	if err != nil {
		return nil, err
	}
	for r.pos < int(size)-1 {
		field, err := r.ReadField()
		if err != nil {
			return nil, err
		}
		_, err = strconv.Atoi(field.Name)
		if err != nil {
			return nil, fmt.Errorf("field %s is not an integer %s", field.Name, err)
		}
		rawVal, err := r.ReadValue(field.Type)
		if err != nil {
			return nil, err
		}
		raw = append(raw, rawVal)
	}
	//fmt.Println("raw", raw)

	return &raw, nil
}

func (r *Reader) ReadSize() (int32, error) {
	var size int32
	err := binary.Read(bytes.NewReader(r.data[r.pos:r.pos+4]), binary.LittleEndian, &size)
	if err != nil {
		return 0, err
	}
	r.pos += 4
	return size, err
}

func (r *Reader) ReadField() (BSONField, error) {
	t := Type(r.data[r.pos])
	r.pos++
	start := r.pos
	for r.data[r.pos] != byte(0) {
		r.pos++
	}
	field := string(r.data[start:r.pos])
	bf := BSONField{t, field}
	r.pos++
	return bf, nil
}

func (r *Reader) ReadValue(t Type) (*Raw, error) {
	switch t {
	case Object:
		raw, err := r.ReadDocValue()
		if err != nil {
			return nil, err
		}
		return raw, nil
	case Double:
		start := r.pos
		r.pos += 8
		return &Raw{Double, r.data[start:r.pos]}, nil
	case String:
		raw, err := r.ReadString()
		if err != nil {
			return nil, err
		}
		return raw, nil
	case Int:
		start := r.pos
		r.pos += 4
		return &Raw{Int, r.data[start:r.pos]}, nil
	case Long:
		start := r.pos
		r.pos += 8
		return &Raw{Long, r.data[start:r.pos]}, nil
	case Array:
		raw, err := r.ReadArrayValue()
		if err != nil {
			return nil, err
		}
		return raw, nil
	case Bool:
		start := r.pos
		r.pos += 1
		return &Raw{Bool, r.data[start:r.pos]}, nil
	}
	return nil, nil
}

func (r *Reader) ReadString() (*Raw, error) {
	length, err := r.ReadSize()
	if err != nil {
		return nil, err
	}
	start := int32(r.pos)
	if r.data[start+length-1] != byte(0) {
		return nil, fmt.Errorf("string length mismatch")
	}
	r.pos = int(start + length)
	return &Raw{String, r.data[start-4 : r.pos]}, nil
}

func (r *Reader) ReadDocValue() (*Raw, error) {
	start, err := r.ReadSized()
	if err != nil {
		return nil, err
	}
	return &Raw{Object, r.data[start-4 : r.pos]}, nil
}

func (r *Reader) ReadArrayValue() (*Raw, error) {
	start, err := r.ReadSized()
	if err != nil {
		return nil, err
	}
	return &Raw{Array, r.data[start-4 : r.pos]}, nil
}

func (r *Reader) ReadSized() (int32, error) {
	length, err := r.ReadSize()
	if err != nil {
		return 0, err
	}
	start := int32(r.pos)
	docBytes := r.data[start : start+length-4]
	if docBytes[len(docBytes)-1] != byte(0) {
		return 0, fmt.Errorf("doc length mismatch")
	}
	r.pos = int(start + length - 4)
	return start, nil
}
