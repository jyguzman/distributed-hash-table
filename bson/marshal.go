package bson

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"reflect"
	"strconv"
)

type Marshaler interface {
	MarshalBSON() ([]byte, error)
}

type ValueMarshaler interface {
	MarshalBSONValue() (Type, []byte, error)
}

type IntValueMarshaler interface {
	ValueMarshaler
	Type() reflect.Type
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
	err := binary.Write(buf, binary.LittleEndian, int32(len(strBytes)+1)) // +1 for null terminator
	err = binary.Write(buf, binary.LittleEndian, strBytes)
	err = binary.Write(buf, binary.LittleEndian, byte(0x00))
	if err != nil {
		return 0, nil, err
	}
	return String, buf.Bytes(), nil
}

func (d D) MarshalBSONValue() (Type, []byte, error) {
	innerBuf := new(bytes.Buffer)
	for _, pair := range d {
		pairBytes, err := pair.MarshalBSON()
		if err != nil {
			return 0, nil, err
		}
		err = binary.Write(innerBuf, binary.LittleEndian, pairBytes)
		if err != nil {
			return 0, nil, err
		}
	}
	buf := new(bytes.Buffer)
	innerBytes := innerBuf.Bytes()
	size := int32(4 + len(innerBytes) + 1) // 4 + 1 to account for the size value and null terminator
	err := binary.Write(buf, binary.LittleEndian, size)
	err = binary.Write(buf, binary.LittleEndian, innerBytes)
	err = binary.Write(buf, binary.LittleEndian, byte(0x00))
	if err != nil {
		return 0, nil, err
	}
	return Object, buf.Bytes(), nil
}

func (m M) MarshalBSONValue() (Type, []byte, error) {
	innerBuf := new(bytes.Buffer)
	for field, val := range m {
		pairBytes, err := Pair{field, val}.MarshalBSON()
		if err != nil {
			return 0, nil, err
		}
		err = binary.Write(innerBuf, binary.LittleEndian, pairBytes)
		if err != nil {
			return 0, nil, err
		}
	}
	buf := new(bytes.Buffer)
	innerBytes := innerBuf.Bytes()
	size := int32(4 + len(innerBytes) + 1) // 4 + 1 to account for the size value and null terminator
	err := binary.Write(buf, binary.LittleEndian, size)
	err = binary.Write(buf, binary.LittleEndian, innerBytes)
	err = binary.Write(buf, binary.LittleEndian, byte(0x00))
	if err != nil {
		return 0, nil, err
	}
	return Object, buf.Bytes(), nil
}

func (a A) MarshalBSONValue() (Type, []byte, error) {
	innerBuf := new(bytes.Buffer)
	for idx, val := range a {
		pairBytes, err := Pair{strconv.Itoa(idx), val}.MarshalBSON()
		if err != nil {
			return 0, nil, err
		}
		err = binary.Write(innerBuf, binary.LittleEndian, pairBytes)
		if err != nil {
			return 0, nil, err
		}
	}
	buf := new(bytes.Buffer)
	innerBytes := innerBuf.Bytes()
	size := int32(4 + len(innerBytes) + 1) // 4 + 1 to account for the size value and null terminator
	err := binary.Write(buf, binary.LittleEndian, size)
	err = binary.Write(buf, binary.LittleEndian, innerBytes)
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

func marshalInt(i any) (Type, []byte, error) {
	buf := new(bytes.Buffer)
	var err error
	var intType Type
	switch t := i.(type) {
	case uint8:
		err = binary.Write(buf, binary.LittleEndian, int32(t))
		intType = Int
	case uint16:
		err = binary.Write(buf, binary.LittleEndian, int32(t))
		intType = Int
	case int8:
		err = binary.Write(buf, binary.LittleEndian, int32(t))
		intType = Int
	case int16:
		err = binary.Write(buf, binary.LittleEndian, int32(t))
		intType = Int
	case int32:
		err = binary.Write(buf, binary.LittleEndian, t)
		intType = Int
	case uint:
		if t <= math.MaxInt32 {
			err = binary.Write(buf, binary.LittleEndian, int32(t))
			intType = Int
		} else {
			err = binary.Write(buf, binary.LittleEndian, int64(t))
			intType = Long
		}
	case uint32:
		if t <= math.MaxInt32 {
			err = binary.Write(buf, binary.LittleEndian, int32(t))
			intType = Int
		} else {
			err = binary.Write(buf, binary.LittleEndian, int64(t))
			intType = Long
		}
	case uint64:
		if t <= math.MaxInt32 {
			err = binary.Write(buf, binary.LittleEndian, int32(t))
			intType = Int
		} else {
			err = binary.Write(buf, binary.LittleEndian, int64(t))
			intType = Long
		}
	case int:
		if t >= math.MinInt32 && t <= math.MaxInt32 {
			err = binary.Write(buf, binary.LittleEndian, int32(t))
			intType = Int
		} else {
			err = binary.Write(buf, binary.LittleEndian, int64(t))
			intType = Long
		}
	case int64:
		err = binary.Write(buf, binary.LittleEndian, t)
		intType = Long
	}
	return intType, buf.Bytes(), err
}

func (bi BSONInt) MarshalBSONValue() (Type, []byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, bi)
	if err != nil {
		return 0, nil, err
	}
	return Int, buf.Bytes(), nil
}

func (bl BSONLong) MarshalBSONValue() (Type, []byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, bl)
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
	case Marshaler:
		data, err := o.MarshalBSON()
		if err != nil {
			return 0, nil, err
		}
		switch o.(type) {
		case A:
			return Array, data, nil
		default:
			return Object, data, nil
		}
	case ValueMarshaler:
		return o.MarshalBSONValue()
	case uint8, uint16, uint32, int8, int16, int32, int64, int, uint:
		return marshalInt(o)
	case float64:
		return BSONDouble(o).MarshalBSONValue()
	case float32:
		return BSONDouble(o).MarshalBSONValue()
	case string:
		return BSONString(o).MarshalBSONValue()
	case []byte:
		return BSONBinData(o).MarshalBSONValue()
	case bool:
		return BSONBool(o).MarshalBSONValue()
	case nil:
		return Null, nil, nil
	default:
		t := reflect.TypeOf(v)
		switch t.Kind() {
		case reflect.Struct:
			return marshalStruct(v)
		case reflect.Slice, reflect.Array:
			val := reflect.ValueOf(v)
			a := A(make([]any, val.Len()))
			for i := 0; i < val.Len(); i++ {
				a[i] = val.Index(i).Interface()
			}
			return MarshalValue(a)
		case reflect.Map:
			val := reflect.ValueOf(v)
			m := M{}
			for _, key := range val.MapKeys() {
				m[key.Interface().(string)] = val.MapIndex(key).Interface()
			}
			return MarshalValue(m)
		case reflect.Ptr:
			val := reflect.Indirect(reflect.ValueOf(v))
			return MarshalValue(val.Interface())
		default:
			return 0, nil, fmt.Errorf("cannot marshal value of type %T", v)
		}
	}
}

func Marshal(v any) ([]byte, error) {
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
		t := reflect.TypeOf(v)
		switch t.Kind() {
		case reflect.Struct:
			_, data, err := MarshalValue(v)
			if err != nil {
				return nil, err
			}
			return data, nil
		case reflect.Slice, reflect.Array:
			val := reflect.ValueOf(v)
			a := A(make([]any, val.Len()))
			for i := 0; i < val.Len(); i++ {
				a[i] = val.Index(i).Interface()
			}
			return Marshal(a)
		case reflect.Map:
			val := reflect.ValueOf(v)
			m := M{}
			for _, key := range val.MapKeys() {
				m[key.Interface().(string)] = val.MapIndex(key).Interface()
			}
			return Marshal(m)
		case reflect.Ptr:
			val := reflect.Indirect(reflect.ValueOf(v))
			return Marshal(val.Interface())
		default:
			return nil, fmt.Errorf("cannot marshal object of type %T", v)
		}
	}
}

func marshalStruct(s any) (Type, []byte, error) {
	rValue := reflect.ValueOf(s)
	rType := rValue.Type()
	if rValue.Kind() != reflect.Struct {
		return 0, nil, fmt.Errorf("value must be a struct")
	}
	innerBuf := new(bytes.Buffer)
	for i := 0; i < rType.NumField(); i++ {
		fieldName := rType.Field(i).Name
		fieldType := rValue.Field(i).Type()
		fieldValue := rValue.Field(i).Interface()
		fieldValType := reflect.TypeOf(fieldValue)
		if fieldType.Kind() == reflect.Interface {
			TypeRegistry[rType.Name()+"."+fieldName] = fieldValType
		}
		pairBytes, err := Pair{Key: fieldName, Val: fieldValue}.MarshalBSON()
		if err != nil {
			return 0, nil, err
		}
		err = binary.Write(innerBuf, binary.LittleEndian, pairBytes)
		if err != nil {
			return 0, nil, err
		}
	}
	buf := new(bytes.Buffer)
	innerBytes := innerBuf.Bytes()
	size := int32(4 + 1 + len(innerBytes))
	err := binary.Write(buf, binary.LittleEndian, size)
	err = binary.Write(buf, binary.LittleEndian, innerBytes)
	err = binary.Write(buf, binary.LittleEndian, byte(0x00))
	if err != nil {
		return 0, nil, err
	}
	return Object, buf.Bytes(), nil
}
