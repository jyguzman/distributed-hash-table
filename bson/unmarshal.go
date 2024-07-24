package bson

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"reflect"
)

type Unmarshaler interface {
	UnmarshalBSON([]byte) error
}

type ValueUnmarshaler interface {
	UnmarshalBSONValue([]byte) error
}

func (rv Raw) Unmarshal(v any) {
	switch rv.Type {
	case Double:
	case String:
	case Int:
	case Long:
	case Bool:
	}
}

func (m M) UnmarshalBSON(b []byte) error {
	r := NewReader(b)
	raw, err := r.ReadDocument()
	if err != nil {
		return err
	}
	for field, val := range raw.Pairs {
		var v any
		err = UnmarshalValue(val.Type, val.Data, &v)
		if err != nil {
			return err
		}
		m[field] = v
	}
	return nil
}

func UnmarshalValue(t Type, v []byte, o *any) error {
	r := bytes.NewReader(v)
	var err error
	switch t {
	case Object, Array:
		m := &M{}
		err = m.UnmarshalBSON(v)
		*o = *m
	case Double:
		var f float64
		err = binary.Read(r, binary.LittleEndian, &f)
		*o = f
	case String:
		*o = string(v)
	case Int:
		//x := o.(*int32)
		err = binary.Read(r, binary.LittleEndian, o)
	case Long:
		//x := o.(*int64)
		err = binary.Read(r, binary.LittleEndian, o)
	case BinData:
		copy((*o).([]byte), v)
	case Bool:
		if v[0] == byte(0) {
			*o = false
		} else {
			*o = true
		}
	default:
		return fmt.Errorf("cannot unmarshal %v into Go value of type %v", t, reflect.TypeOf(o))
	}
	if err != nil {
		return err
	}
	return nil
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
