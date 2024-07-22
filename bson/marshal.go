package bson

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/big"
	"reflect"
	"strconv"
)

type Marshaler interface {
	MarshalBSON() ([]byte, error)
}

type ValueMarshaler interface {
	MarshalBSONValue() (Type, []byte, error)
}

func (bd BSONDouble) MarshalBSONValue() (Type, []byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, float64(bd))
	if err != nil {
		return 0, nil, err
	}
	return Double, buf.Bytes(), nil
}

func (bs BSONString) MarshalBSONValue() (Type, []byte, error) {
	buf, strBytes := new(bytes.Buffer), []byte(bs)
	err := binary.Write(buf, binary.LittleEndian, int32(len(strBytes)))
	err = binary.Write(buf, binary.LittleEndian, strBytes)
	err = binary.Write(buf, binary.LittleEndian, int32(0x00))
	if err != nil {
		return 0, nil, err
	}
	return String, buf.Bytes(), nil
}

func (d D) MarshalBSONValue() (Type, []byte, error) {
	size := int32(4) // starts at 5 because of the size value itself and null terminator
	var pairs []byte
	for _, pair := range d {
		pairBytes, err := pair.MarshalBSON()
		if err != nil {
			return 0, nil, err
		}
		size += int32(len(pairBytes))
		pairs = append(pairs, pairBytes...)
	}
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, size)
	err = binary.Write(buf, binary.LittleEndian, pairs)
	err = binary.Write(buf, binary.LittleEndian, byte(0x00))
	if err != nil {
		return 0, nil, err
	}
	return Object, buf.Bytes(), nil
}

func (m M) MarshalBSONValue() (Type, []byte, error) {
	size := int32(4) // starts at 5 because of the size value itself and null terminator
	var pairs []byte
	for field, val := range m {
		pairBytes, err := Pair{field, val}.MarshalBSON()
		if err != nil {
			return 0, nil, err
		}
		size += int32(len(pairBytes))
		pairs = append(pairs, pairBytes...)
	}
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, size)
	err = binary.Write(buf, binary.LittleEndian, pairs)
	err = binary.Write(buf, binary.LittleEndian, byte(0x00))
	if err != nil {
		return 0, nil, err
	}
	return Object, buf.Bytes(), nil
}

func (a A) MarshalBSONValue() (Type, []byte, error) {
	size := int32(5) // starts at 5 because of the size value itself and null terminator
	var pairs []byte
	for idx, val := range a {
		pairBytes, err := Pair{strconv.Itoa(idx), val}.MarshalBSON()
		if err != nil {
			return 0, nil, err
		}
		size += int32(len(pairBytes))
		pairs = append(pairs, pairBytes...)
	}
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, size)
	err = binary.Write(buf, binary.LittleEndian, pairs)
	err = binary.Write(buf, binary.LittleEndian, byte(0x00))
	if err != nil {
		return 0, nil, err
	}
	return Array, buf.Bytes(), nil
}

func (bd BSONBinData) MarshalBSONValue() (Type, []byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, int32(len(bd)))
	err = binary.Write(buf, binary.LittleEndian, bd)
	if err != nil {
		return 0, nil, err
	}
	return BinData, buf.Bytes(), nil
}

func (bb BSONBool) MarshalBSONValue() (Type, []byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, bb)
	if err != nil {
		return 0, nil, err
	}
	return Bool, buf.Bytes(), nil
}

func (bi BSONInt) MarshalBSONValue() (Type, []byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, int32(bi))
	if err != nil {
		return 0, nil, err
	}
	return Int, buf.Bytes(), nil
}

func (bl BSONLong) MarshalBSONValue() (Type, []byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, int64(bl))
	if err != nil {
		return 0, nil, err
	}
	return Long, buf.Bytes(), nil
}

func (bf BSONField) MarshalBSON() ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, bf.Type)
	err = binary.Write(buf, binary.LittleEndian, []byte(bf.Name))
	err = binary.Write(buf, binary.LittleEndian, byte(0x00))
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (d D) MarshalBSON() ([]byte, error) {
	_, data, err := d.MarshalBSONValue()
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (m M) MarshalBSON() ([]byte, error) {
	_, data, err := m.MarshalBSONValue()
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (a A) MarshalBSON() ([]byte, error) {
	_, data, err := a.MarshalBSONValue()
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (p Pair) MarshalBSON() ([]byte, error) {
	buf := new(bytes.Buffer)
	field, val := p.Key, p.Val
	valType, valBytes, err := MarshalValue(val)
	if err != nil {
		return nil, err
	}
	bfBytes, err := BSONField{valType, field}.MarshalBSON()
	if err != nil {
		return nil, err
	}
	err = binary.Write(buf, binary.LittleEndian, bfBytes)
	if err != nil {
		return nil, err
	}
	err = binary.Write(buf, binary.LittleEndian, valBytes)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func MarshalValue(v any) (Type, []byte, error) {
	switch o := v.(type) {
	case ValueMarshaler:
		return o.MarshalBSONValue()
	case string:
		return BSONString(o).MarshalBSONValue()
	case []byte:
		return BSONBinData(o).MarshalBSONValue()
	case int64:
		return BSONLong(o).MarshalBSONValue()
	case int32:
		return BSONInt(o).MarshalBSONValue()
	case bool:
		return BSONBool(o).MarshalBSONValue()
	case float64:
		return BSONDouble(o).MarshalBSONValue()
	default:
		return 0, nil, fmt.Errorf("cannot marshal value of type %T", v)
	}
}

func marshal(v any) ([]byte, error) {
	switch o := v.(type) {
	case Marshaler:
		return o.MarshalBSON()
	case ValueMarshaler:
		_, data, err := o.MarshalBSONValue()
		if err != nil {
			return nil, err
		}
		return data, nil
	default:
		return nil, fmt.Errorf("cannot marshal value of type %T", v)
	}
}

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

func marshalValue(val any) (Type, []byte, error) {
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
	valType, valBytes, err := marshalValue(p.Val)
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
