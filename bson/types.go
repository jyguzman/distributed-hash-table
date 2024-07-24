package bson

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

type BSONField struct {
	Type Type
	Name string
}
type BSONString string
type BSONBool bool
type BSONDouble float64
type BSONInt int32
type BSONLong int64
type BSONBinData []byte

type Raw struct {
	Type Type
	Data []byte
}

type RawD struct {
	Size   int32
	Fields []BSONField
	Values []*Raw
	Pairs  map[string]*Raw
}
