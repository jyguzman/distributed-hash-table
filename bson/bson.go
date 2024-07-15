package bson

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/big"
	"slices"
	"strconv"
)

type Type int8

const (
	Double  Type = 0x01
	String  Type = 0x02
	Object  Type = 0x03
	Array   Type = 0x04
	BinData Type = 0x05
	Bool    Type = 0x08
	Null    Type = 0xA
	Int     Type = 0x10
	Long    Type = 0x12
)

type Pair struct {
	Key string
	Val any
}

type D []Pair
type M map[string]any
type A []any

func Marshal(v any) ([]byte, error) {
	buf := new(bytes.Buffer)
	var err error
	switch t := v.(type) {
	case D:
		objBytes, size, marshalErr := marshalObj(t)
		if marshalErr != nil {
			return nil, marshalErr
		}
		size += 4 + 1 // the size of the size int32 value itself plus the null terminator
		err = binary.Write(buf, binary.LittleEndian, size)
		err = binary.Write(buf, binary.LittleEndian, objBytes)
	case M:
		mBytes, size, marshalErr := marshalMap(t)
		if marshalErr != nil {
			return nil, marshalErr
		}
		size += 4 + 1
		err = binary.Write(buf, binary.LittleEndian, size)
		err = binary.Write(buf, binary.LittleEndian, mBytes)
	case A:
		aBytes, size, marshalErr := marshalArray(t)
		if marshalErr != nil {
			return nil, marshalErr
		}
		size += 4 + 1
		err = binary.Write(buf, binary.LittleEndian, size)
		err = binary.Write(buf, binary.LittleEndian, aBytes)
	default:
		return nil, fmt.Errorf("value must be object (ordered/unordered) or array, not of type %T", t)
	}
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func marshalMap(m M) ([]byte, int32, error) {
	buf := new(bytes.Buffer)
	var size int32
	for key, val := range m {
		pairBytes, pairSize, err := MarshalPair(Pair{Key: key, Val: val})
		if err != nil {
			return nil, -1, err
		}
		size += pairSize
		_, err = buf.Write(pairBytes)
	}
	err := binary.Write(buf, binary.LittleEndian, byte(0x00))
	if err != nil {
		return nil, -1, err
	}
	return buf.Bytes(), size, nil
}

func marshalObj(obj D) ([]byte, int32, error) {
	buf := new(bytes.Buffer)
	var size int32
	for _, pair := range obj {
		pairBytes, pairSize, err := MarshalPair(pair)
		if err != nil {
			return nil, -1, err
		}
		size += pairSize
		_, err = buf.Write(pairBytes)
	}
	err := binary.Write(buf, binary.LittleEndian, byte(0x00))
	if err != nil {
		return nil, -1, err
	}
	return buf.Bytes(), size, nil
}

func marshalArray(array A) ([]byte, int32, error) {
	buf := new(bytes.Buffer)
	var size int32
	for idx, val := range array {
		pairBytes, pairSize, err := MarshalPair(Pair{Key: strconv.Itoa(idx), Val: val})
		if err != nil {
			return nil, -1, err
		}
		size += pairSize
		_, err = buf.Write(pairBytes)
	}
	err := binary.Write(buf, binary.LittleEndian, byte(0x00))
	if err != nil {
		return nil, -1, err
	}
	return buf.Bytes(), size, nil
}

func MarshalValue(val any) (Type, []byte, error) {
	buf := new(bytes.Buffer)
	var bsonType Type
	var err error
	if val == nil {
		err = binary.Write(buf, binary.LittleEndian, int8(0x00))
		if err != nil {
			return -1, nil, err
		}
		return Null, buf.Bytes(), nil
	}
	switch vt := val.(type) {
	case D:
		bsonType = Object
		objBytes, marshalErr := Marshal(vt)
		if marshalErr != nil {
			return -1, nil, marshalErr
		}
		err = binary.Write(buf, binary.LittleEndian, objBytes)
	case M:
		bsonType = Object
		mapBytes, marshalErr := Marshal(vt)
		if marshalErr != nil {
			return -1, nil, marshalErr
		}
		err = binary.Write(buf, binary.LittleEndian, mapBytes)
	case A:
		bsonType = Array
		arrBytes, marshalErr := Marshal(vt)
		if marshalErr != nil {
			return -1, nil, marshalErr
		}
		err = binary.Write(buf, binary.LittleEndian, arrBytes)
	case []byte:
		bsonType = BinData
		err = binary.Write(buf, binary.LittleEndian, vt)
	case float64:
		bsonType = Double
		err = binary.Write(buf, binary.LittleEndian, vt)
	case int32:
		bsonType = Int
		err = binary.Write(buf, binary.LittleEndian, vt)
	case int:
		bsonType = Long
		err = binary.Write(buf, binary.LittleEndian, int64(vt))
	case int64:
		bsonType = Long
		err = binary.Write(buf, binary.LittleEndian, vt)
	case *big.Int:
		bsonType = BinData
		intBytes := vt.Bytes()
		slices.Reverse(intBytes)
		err = binary.Write(buf, binary.LittleEndian, intBytes)
	case string:
		bsonType = String
		err = binary.Write(buf, binary.LittleEndian, []byte(vt))
	case bool:
		bsonType = Bool
		err = binary.Write(buf, binary.LittleEndian, vt)
	default:
		return -1, nil, fmt.Errorf("unsupported type: %T", val)
	}
	if err != nil {
		return -1, nil, fmt.Errorf("error serializing value: %v", val)
	}
	return bsonType, buf.Bytes(), err
}

func MarshalPair(p Pair) ([]byte, int32, error) {
	buf := new(bytes.Buffer)
	valType, valBytes, err := MarshalValue(p.Val)
	if err != nil {
		return nil, -1, err
	}
	err = binary.Write(buf, binary.LittleEndian, valType)
	err = binary.Write(buf, binary.LittleEndian, []byte(p.Key))
	err = binary.Write(buf, binary.LittleEndian, byte(0x00))
	if valType == String {
		err = binary.Write(buf, binary.LittleEndian, int32(len(valBytes)))
	}
	err = binary.Write(buf, binary.LittleEndian, valBytes)
	if valType == String {
		err = binary.Write(buf, binary.LittleEndian, byte(0x00))
	}
	if err != nil {
		return nil, -1, err
	}
	bufBytes := buf.Bytes()
	return bufBytes, int32(len(bufBytes)), nil
}

func Unmarshal(data []byte, v any) {

}
