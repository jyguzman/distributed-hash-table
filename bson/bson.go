package bson

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math/big"
	"reflect"
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
		if reflect.TypeOf(v).Kind() == reflect.Struct {
			sBytes, size, marshalErr := marshalStruct(v)
			if marshalErr != nil {
				return nil, marshalErr
			}
			size += 4 + 1
			err = binary.Write(buf, binary.LittleEndian, size)
			err = binary.Write(buf, binary.LittleEndian, sBytes)
		} else {
			return nil, fmt.Errorf("value must be object (ordered/unordered) or array, not of type %T", t)
		}
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

func marshalStruct(s any) ([]byte, int32, error) {
	rValue := reflect.ValueOf(s)
	rType := rValue.Type()
	if rValue.Kind() != reflect.Struct {
		return nil, 0, fmt.Errorf("value must be a struct")
	}
	buf := new(bytes.Buffer)
	var size int32
	for i := 0; i < rType.NumField(); i++ {
		fieldName := rType.Field(i).Name
		fieldValue := rValue.Field(i).Interface()
		pairBytes, pairSize, err := MarshalPair(Pair{Key: fieldName, Val: fieldValue})
		if err != nil {
			return nil, 0, err
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
	case M, *M:
		bsonType = Object
		mapBytes, marshalErr := Marshal(vt)
		if marshalErr != nil {
			return -1, nil, marshalErr
		}
		err = binary.Write(buf, binary.LittleEndian, mapBytes)
	case A, []A:
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
		bsonType = String
		hexString := vt.Text(16)
		err = binary.Write(buf, binary.LittleEndian, []byte(hexString))
	case string:
		bsonType = String
		err = binary.Write(buf, binary.LittleEndian, []byte(vt))
	case bool:
		bsonType = Bool
		err = binary.Write(buf, binary.LittleEndian, vt)
	default:
		if reflect.TypeOf(val).Kind() == reflect.Struct {
			bsonType = Object
			sBytes, marshalErr := Marshal(vt)
			if marshalErr != nil {
				return -1, nil, marshalErr
			}
			err = binary.Write(buf, binary.LittleEndian, sBytes)
		} else {
			return -1, nil, fmt.Errorf("unsupported type: %T", val)
		}
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
		err = binary.Write(buf, binary.LittleEndian, int32(len(valBytes)+1))
	}
	if valType == BinData {
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

func Unmarshal(data []byte, obj any) (uint32, error) {
	rValue := reflect.ValueOf(obj)
	rType := rValue.Type()
	if rType.Kind() == reflect.Ptr && rValue.Elem().Kind() == reflect.Struct {
		return UnmarshalStruct(data, obj)
	}
	size := binary.LittleEndian.Uint32(data[:4])
	if data[size-1] != byte(0x00) {
		return 0, fmt.Errorf("last byte must be null terminator 0x00, is 0x%02x", data[size-1])
	}
	i := uint32(4)
	for i < size-1 {
		bsonType := Type(data[i])
		i++
		fieldStart := i
		for data[i] != byte(0x00) {
			i++
		}
		field := string(data[fieldStart:i])
		i++

		if i > size-1 {
			break
		}

		var (
			value any
			skip  uint32
			err   error
		)

		if bsonType == Object || bsonType == Array {
			switch obj.(type) {
			case *D, A:
				inner := D{}
				skip, err = Unmarshal(data[i:], &inner)
				if err != nil {
					return 0, err
				}
				value = Pair{Key: field, Val: inner}
			case M, *M:
				value = M{}
				skip, err = Unmarshal(data[i:], value)
				if err != nil {
					return 0, err
				}
			default:
				return 0, fmt.Errorf("not an object or array type: %T", obj)
			}
		} else {
			value, skip, err = unmarshalValue(data[i:], bsonType)
			if err != nil {
				return 0, err
			}
		}
		i += skip
		switch ot := obj.(type) {
		case M:
			ot[field] = value
		case *M:
			(*ot)[field] = value
		case *D:
			switch value.(type) {
			case Pair:
				*ot = append(*ot, value.(Pair))
			default:
				*ot = append(*ot, Pair{Key: field, Val: value})
			}
		}
		if data[i] == 0x00 {
			i++
		}
	}
	return i, nil
}

func UnmarshalStruct(data []byte, obj any) (uint32, error) {
	rValue := reflect.ValueOf(obj).Elem()
	if rValue.Kind() != reflect.Struct {
		return 0, fmt.Errorf("expected struct, got %T", obj)
	}
	m := M{}
	i, err := Unmarshal(data, &m)
	if err != nil {
		return 0, err
	}
	mBytes, err := json.Marshal(m)
	if err != nil {
		return 0, err
	}
	err = json.Unmarshal(mBytes, obj)
	if err != nil {
		return 0, err
	}
	return i, nil
}

func unmarshalValue(v []byte, vType Type) (any, uint32, error) {
	switch vType {
	case String:
		length := binary.LittleEndian.Uint32(v[:4])
		end := length + 3
		if v[end] != 0x00 {
			return nil, 0, fmt.Errorf("strings must have null terminator, has 0x%02x", v[end])
		}
		return string(v[4:end]), end + 1, nil
	case Int:
		return int32(binary.LittleEndian.Uint32(v[:4])), 4, nil
	case Long:
		return int64(binary.LittleEndian.Uint64(v[:8])), 8, nil
	case Double:
		var float float64
		buf := bytes.NewReader(v[:8])
		err := binary.Read(buf, binary.LittleEndian, &float)
		if err != nil {
			return nil, 0, err
		}
		return float, 8, nil
	case Bool:
		if v[0] == 0x00 {
			return false, 1, nil
		}
		return true, 1, nil
	case BinData:
		length := binary.LittleEndian.Uint32(v[:4])
		end := length + 4
		return v[4:end], end, nil
	case Null:
		return nil, 1, nil
	}
	return nil, 0, nil
}
