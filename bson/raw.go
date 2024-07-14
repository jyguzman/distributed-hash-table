package bson

import "fmt"

type Raw []byte

func (r Raw) Size() int32 { return int32(len(r)) }

func (r Raw) Lookup(fields ...string) any {
	start, end := 0, len(fields)-1

	fmt.Println(start, end)
	return nil
}
