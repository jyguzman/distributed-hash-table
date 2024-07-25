package bson

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"reflect"
)

type Unmarshaler interface {
	UnmarshalBSON([]byte) error
}

type ValueUnmarshaler interface {
	UnmarshalBSONValue(Type, []byte) error
}

func (rv Raw) Unmarshal(v any) error {
	if v == nil || reflect.TypeOf(v).Kind() != reflect.Ptr {
		return fmt.Errorf("value must be non-nil and a pointer")
	}
	var err error
	switch rv.Type {
	case Double:
		t, ok := v.(*float64)
		if !ok {
			return fmt.Errorf("cannot unmarshal Double into %T", t)
		}
		err = binary.Read(bytes.NewReader(rv.Data), binary.LittleEndian, t)
	case String:
		t, ok := v.(*string)
		if !ok {
			return fmt.Errorf("cannot unmarshal String into %T", t)
		}
		*t = string(rv.Data)[4 : len(rv.Data)-1]
	case Int:
		t, ok := v.(*int32)
		if !ok {
			return fmt.Errorf("cannot unmarshal Int into %T", t)
		}
		err = binary.Read(bytes.NewReader(rv.Data), binary.LittleEndian, t)
	case Long:
		t, ok := v.(*int64)
		if !ok {
			return fmt.Errorf("cannot unmarshal Long into %T", t)
		}
		err = binary.Read(bytes.NewReader(rv.Data), binary.LittleEndian, t)
	case Bool:
		t, ok := v.(*bool)
		if !ok {
			return fmt.Errorf("cannot unmarshal Bool into %T", t)
		}
		err = binary.Read(bytes.NewReader(rv.Data), binary.LittleEndian, t)
	case Array:
		t, ok := v.(*A)
		if !ok {
			return fmt.Errorf("cannot unmarshal Array into %T", t)
		}
		err = t.UnmarshalBSON(rv.Data)
	case Object:
		t, ok := v.(Unmarshaler)
		if !ok {
			return fmt.Errorf("cannot unmarshal Object into %T", t)
		}
		err = t.UnmarshalBSON(rv.Data)
	}
	return err
}

func (d *D) UnmarshalBSON(b []byte) error {
	r := NewReader(b)
	raw, err := r.ReadDocument()
	if err != nil {
		return err
	}

	for field, val := range raw.Pairs {
		var value any
		switch val.Type {
		case Double:
			v := new(float64)
			err = UnmarshalValue(val.Type, val.Data, v)
			value = *v
		case String:
			v := new(string)
			err = UnmarshalValue(val.Type, val.Data, v)
			value = *v
		case Int:
			v := new(int32)
			err = UnmarshalValue(val.Type, val.Data, v)
			value = *v
		case Long:
			v := new(int64)
			err = UnmarshalValue(val.Type, val.Data, v)
			value = *v
		case Bool:
			v := new(bool)
			err = UnmarshalValue(val.Type, val.Data, v)
			value = *v
		case Object:
			v := new(D)
			err = UnmarshalValue(val.Type, val.Data, v)
			value = *v
		case Null:
		case Array:
			v := new(A)
			err = v.UnmarshalWithParent(val.Data, d)
			value = *v
		}
		*d = append(*d, Pair{Key: field, Val: value})
	}
	return nil
}

func (m *M) UnmarshalBSON(b []byte) error {
	r := NewReader(b)
	raw, err := r.ReadDocument()
	if err != nil {
		return err
	}

	for field, val := range raw.Pairs {
		//fmt.Println("field:", field, "val:", val)
		var value any
		switch val.Type {
		case Double:
			v := new(float64)
			err = UnmarshalValue(val.Type, val.Data, v)
			value = *v
		case String:
			v := new(string)
			err = UnmarshalValue(val.Type, val.Data, v)
			value = *v
		case Int:
			v := new(int32)
			err = UnmarshalValue(val.Type, val.Data, v)
			value = *v
		case Long:
			v := new(int64)
			err = UnmarshalValue(val.Type, val.Data, v)
			value = *v
		case Bool:
			v := new(bool)
			err = UnmarshalValue(val.Type, val.Data, v)
			value = *v
		case Object:
			v := M{}
			err = UnmarshalValue(val.Type, val.Data, &v)
			value = v
		case Null:
		case Array:
			v := new(A)
			err = v.UnmarshalBSON(val.Data)
			//fmt.Println("v:", v)
			value = *v
		}
		//fmt.Println("value:", value)
		(*m)[field] = value
	}
	return nil
}

func (a *A) UnmarshalWithParent(b []byte, o any) error {
	r := NewReader(b)
	raw, err := r.ReadArray()
	if err != nil {
		return err
	}

	for _, val := range *raw {
		var value any
		switch val.Type {
		case Double:
			v := new(float64)
			err = UnmarshalValue(val.Type, val.Data, v)
			value = *v
		case String:
			v := new(string)
			err = UnmarshalValue(val.Type, val.Data, v)
			value = *v
		case Int:
			v := new(int32)
			err = UnmarshalValue(val.Type, val.Data, v)
			value = *v
		case Long:
			v := new(int64)
			err = UnmarshalValue(val.Type, val.Data, v)
			value = *v
		case Bool:
			v := new(bool)
			err = UnmarshalValue(val.Type, val.Data, v)
			value = *v
		case Object:
			switch o.(type) {
			case *D:
				v := new(D)
				err = UnmarshalValue(val.Type, val.Data, v)
				value = *v
			case *M:
				v := M{}
				err = UnmarshalValue(val.Type, val.Data, &v)
				value = v
			default:
				return fmt.Errorf("cannot unmarshal Object into %T", o)
			}
		case Null:
		case Array:
			v := new(A)
			err = UnmarshalValue(val.Type, val.Data, v)
			value = *v
		}
		if err != nil {
			return err
		}
		*a = append(*a, value)
	}
	return nil
}

func (a *A) UnmarshalBSON(b []byte) error {
	r := NewReader(b)
	raw, err := r.ReadArray()
	if err != nil {
		return err
	}

	for _, val := range *raw {
		var value any
		switch val.Type {
		case Double:
			v := new(float64)
			err = UnmarshalValue(val.Type, val.Data, v)
			value = *v
		case String:
			v := new(string)
			err = UnmarshalValue(val.Type, val.Data, v)
			value = *v
		case Int:
			v := new(int32)
			err = UnmarshalValue(val.Type, val.Data, v)
			value = *v
		case Long:
			v := new(int64)
			err = UnmarshalValue(val.Type, val.Data, v)
			value = *v
		case Bool:
			v := new(bool)
			err = UnmarshalValue(val.Type, val.Data, v)
			value = *v
		case Object:
			v := M{}
			err = UnmarshalValue(val.Type, val.Data, &v)
			value = v
		case Null:
		case Array:
			v := new(A)
			err = UnmarshalValue(val.Type, val.Data, v)
			value = *v
		}
		if err != nil {
			return err
		}
		*a = append(*a, value)
	}
	return nil
}

func UnmarshalValue(t Type, v []byte, o any) error {
	return Raw{t, v}.Unmarshal(o)
}

func Unmarshal(data []byte, obj any) error {
	rValue := reflect.ValueOf(obj)
	rType := rValue.Type()
	if rType.Kind() != reflect.Ptr {
		return fmt.Errorf("object to unmarshal into must be a pointer")
	}
	switch t := obj.(type) {
	case nil:
		return fmt.Errorf("cannot unmarshal into nil object")
	case Unmarshaler:
		return t.UnmarshalBSON(data)
	case *float64:
		return UnmarshalValue(Double, data, obj)
	case *string:
		return UnmarshalValue(String, data, obj)
	case *int32:
		return UnmarshalValue(Int, data, obj)
	case *int64:
		return UnmarshalValue(Long, data, obj)
	case *bool:
		return UnmarshalValue(Bool, data, obj)
	default:
		rIndirect := reflect.Indirect(rValue)
		if rIndirect.Kind() == reflect.Struct {
			m := M{}
			err := Unmarshal(data, &m)
			if err != nil {
				return err
			}
			err = UnmarshalStruct(m, obj)
			return nil
		}
		return fmt.Errorf("cannot unmarshal into %T", obj)
	}
}

func UnmarshalStruct(m M, obj any) error {
	rValue := reflect.ValueOf(obj)
	rType := rValue.Type()
	if rType.Kind() != reflect.Ptr {
		return fmt.Errorf("object to unmarshal into must be a pointer")
	}
	rIndirect := reflect.Indirect(rValue)
	if rIndirect.Kind() != reflect.Struct {
		return fmt.Errorf("object to unmarshal into must be a struct")
	}

	for k, v := range m {
		typ := reflect.TypeOf(v)
		switch typ.Kind() {
		case reflect.Map:
			mBytes, err := Marshal(m[k].(M))
			if err != nil {
				return err
			}
			sfType := rValue.Elem().FieldByName(k).Type()
			newS := reflect.New(sfType).Interface()
			err = Unmarshal(mBytes, newS)
			if err != nil {
				return err
			}
			rValue.Elem().FieldByName(k).Set(reflect.ValueOf(newS).Elem())
		case reflect.Array, reflect.Slice:
			aType := rValue.Elem().FieldByName(k).Type()
			arr := v.(A)
			newA := reflect.MakeSlice(aType, len(arr), len(arr))
			elemType := reflect.TypeOf(newA.Index(0).Interface())
			for i := 0; i < newA.Len(); i++ {
				arrElemType := reflect.TypeOf(arr[i])
				if arrElemType.Kind() == reflect.Map {
					arrMBytes, err := Marshal(arr[i].(M))
					newElem := reflect.New(elemType).Interface()
					err = Unmarshal(arrMBytes, newElem)
					if err != nil {
						return err
					}
					newA.Index(i).Set(reflect.Indirect(reflect.ValueOf(newElem)))
				} else {
					newA.Index(i).Set(reflect.ValueOf(arr[i]))
				}
			}
			rValue.Elem().FieldByName(k).Set(newA)
		default:
			rValue.Elem().FieldByName(k).Set(reflect.ValueOf(v))
		}
	}
	return nil
}

func unmarshalArray() error {
	return nil
}
