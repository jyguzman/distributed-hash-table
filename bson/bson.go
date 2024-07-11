package bson

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type Type int8

const (
	Double  Type = 0x01
	Int     Type = 0x02
	Long    Type = 0x03
	String  Type = 0x04
	Boolean Type = 0x05
	Array   Type = 0x06
	Object  Type = 0x07
	Null    Type = 0x08
)

type Pair struct {
	Key string
	Val any
}

type D []Pair
type O map[string]any

func Marshal(v any) ([]byte, error) {
	return nil, nil
}

func MarshalValue(val any) ([]byte, error) {
	buf := new(bytes.Buffer)
	var err error
	if val == nil {
		err = binary.Write(buf, binary.BigEndian, Null)
		err = binary.Write(buf, binary.BigEndian, int8(0x00))
		if err != nil {
			return nil, err
		}
		return buf.Bytes(), nil
	}
	switch vt := val.(type) {
	case float64:
		err = binary.Write(buf, binary.BigEndian, Double)
		err = binary.Write(buf, binary.BigEndian, vt)
	case int32:
		err = binary.Write(buf, binary.BigEndian, Int)
		err = binary.Write(buf, binary.BigEndian, vt)
	case int:
		err = binary.Write(buf, binary.BigEndian, Long)
		err = binary.Write(buf, binary.BigEndian, int64(vt))
	case int64:
		err = binary.Write(buf, binary.BigEndian, Long)
		err = binary.Write(buf, binary.BigEndian, vt)
	case string:
		err = binary.Write(buf, binary.BigEndian, String)
		err = binary.Write(buf, binary.BigEndian, []byte(vt))
	case bool:
		err = binary.Write(buf, binary.BigEndian, Boolean)
		err = binary.Write(buf, binary.BigEndian, vt)
	default:
		return nil, fmt.Errorf("unsupported type: %T", val)
	}
	if err != nil {
		return nil, fmt.Errorf("error serializing value: %v", val)
	}
	return buf.Bytes(), err
}
