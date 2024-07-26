package bson

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"reflect"
	"slices"
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
	case Null:
		v = nil
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
			value = nil
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
			value = nil
		case Array:
			v := new(A)
			err = v.UnmarshalBSON(val.Data)
			value = *v
		}
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
		return UnmarshalValue(Null, data, obj)
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
			err = unmarshalStruct(m, obj)
			return nil
		}
		return fmt.Errorf("cannot unmarshal %v into %T", rValue, obj)
	}
}

func unmarshalStruct(m M, obj any) error {
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
		vType := reflect.TypeOf(v)
		field := rValue.Elem().FieldByName(k)
		fieldType := field.Type()
		var valueToSetType reflect.Type
		if v == nil {
			field.Set(reflect.Zero(field.Type()))
			continue
		}
		if fieldType.Kind() == reflect.Interface {
			regKey := rType.Elem().Name() + "." + k
			typ, exists := TypeRegistry[regKey]
			if !exists {
				valueToSetType = reflect.TypeOf(v)
			} else {
				valueToSetType = typ
				//delete(TypeRegistry, k)
			}

		} else {
			valueToSetType = fieldType
		}
		switch vType.Kind() {
		case reflect.Map:
			newStruct, err := StructFromBSONMap(v.(M), valueToSetType)
			if err != nil {
				return err
			}
			field.Set(reflect.ValueOf(newStruct).Elem())
		case reflect.Array, reflect.Slice:
			newArray, err := sliceFromBSONArray(v.(A), valueToSetType)
			if err != nil {
				return err
			}
			field.Set(*newArray)
		case reflect.Int32, reflect.Int64:
			if field.Type().Kind() == reflect.Interface {
				field.Set(reflect.ValueOf(v))
			} else {
				err := setStructIntegerField(field, v)
				if err != nil {
					return err
				}
			}
		case reflect.Float64:
			if fieldType.Kind() == reflect.Interface {
				field.Set(reflect.ValueOf(v))
			} else {
				field.SetFloat(v.(float64))
			}
		case reflect.String:
			if fieldType.Kind() == reflect.Interface {
				field.Set(reflect.ValueOf(v))
			} else {
				field.SetString(v.(string))
			}
		case reflect.Bool:
			if fieldType.Kind() == reflect.Interface {
				field.Set(reflect.ValueOf(v))
			} else {
				field.SetBool(v.(bool))
			}
		default:
			return fmt.Errorf("cannot unmarshal into %T", m[k])
		}
	}
	return nil
}

func structTypeFromBSONMap(m M) (reflect.Type, error) {
	var structFields []reflect.StructField
	var err error
	for key, val := range m {
		typ := reflect.TypeOf(val)
		if typ.Kind() == reflect.Map {
			typ, err = structTypeFromBSONMap(val.(M))
			if err != nil {
				return nil, err
			}
		}
		structFields = append(structFields, reflect.StructField{
			Name: key,
			Type: typ,
		})
	}
	return reflect.StructOf(structFields), nil
}

func StructFromBSONMap(m M, valueType reflect.Type) (any, error) {
	mBytes, err := Marshal(m)
	if err != nil {
		return nil, err
	}
	var newStruct any
	if valueType.Kind() == reflect.Map {
		var structFields []reflect.StructField
		for key, val := range m {
			valType := reflect.TypeOf(val)
			if reflect.TypeOf(val).Kind() == reflect.Map {
				valType, err = structTypeFromBSONMap(val.(M))
				if err != nil {
					return nil, err
				}
			}
			structFields = append(structFields, reflect.StructField{
				Name: key,
				Type: valType,
			})
		}
		typ := reflect.StructOf(structFields)
		newStruct = reflect.New(typ).Interface()
	} else {
		newStruct = reflect.New(valueType).Interface()
	}
	err = Unmarshal(mBytes, newStruct)
	if err != nil {
		return nil, err
	}
	return newStruct, nil
}

func sliceFromBSONArray(fromArray A, newArrayType reflect.Type) (*reflect.Value, error) {
	newArray := reflect.MakeSlice(newArrayType, len(fromArray), len(fromArray))
	for i := 0; i < newArray.Len(); i++ {
		elemType := reflect.TypeOf(newArray.Index(i).Interface())
		if reflect.TypeOf(fromArray[i]).Kind() == reflect.Map {
			newStruct, err := StructFromBSONMap(fromArray[i].(M), elemType)
			if err != nil {
				return nil, err
			}
			newArray.Index(i).Set(reflect.ValueOf(newStruct).Elem())
		} else {
			newArray.Index(i).Set(reflect.ValueOf(fromArray[i]))
		}
	}
	return &newArray, nil
}

func setStructIntegerField(field reflect.Value, v any) error {
	rValue, vType, fieldType := reflect.ValueOf(v), reflect.TypeOf(v), field.Type()
	intKinds := []reflect.Kind{
		reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int,
		reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
	}
	if !slices.Contains(intKinds, field.Kind()) {
		return fmt.Errorf("cannot set %T into field of type %s", rValue.Kind(), fieldType)
	}
	if !slices.Contains(intKinds, vType.Kind()) {
		return fmt.Errorf("val (%T, %v) is not int32 or int64", v, v)
	}
	switch vType.Kind() {
	case reflect.Int32:
		switch fieldType.Kind() {
		case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
			field.SetInt(int64(v.(int32)))
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			field.SetUint(uint64(v.(int32)))
		default:
			return fmt.Errorf("cannot marshal %T into int", v)
		}
	case reflect.Int64:
		switch fieldType.Kind() {
		case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
			field.SetInt(v.(int64))
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			field.SetUint(uint64(v.(int64)))
		default:
			return fmt.Errorf("cannot marshal %T into int", v)
		}
	default:
		return fmt.Errorf("%T is not an int32 or int64", v)
	}
	return nil
}

func setStructField(field reflect.Value, v any) error {
	if !field.IsValid() {
		return fmt.Errorf("field %v is not valid", field)
	}
	if !field.CanSet() {
		return fmt.Errorf("field %v is not settable", field)
	}
	vType := reflect.TypeOf(v)
	switch vType.Kind() {
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int,
		reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return setStructIntegerField(field, v)
	case reflect.Float32, reflect.Float64:
		field.SetFloat(v.(float64))
	case reflect.String:
		field.SetString(v.(string))
	case reflect.Bool:
		field.SetBool(v.(bool))
	case reflect.Slice, reflect.Array:
		arrayType := field.Type()
		newArray, err := sliceFromBSONArray(v.(A), arrayType)
		if err != nil {
			return err
		}
		field.Set(*newArray)
	case reflect.Map:
		newStruct, err := StructFromBSONMap(v.(M), vType)
		if err != nil {
			return err
		}
		field.Set(reflect.ValueOf(newStruct).Elem())
	default:
		return fmt.Errorf("unhandled field type %v", field.Type())
	}
	return nil
}
