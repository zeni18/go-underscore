package __

import (
	"fmt"
	"reflect"
)

func init() {
	fmt.Println()
	MakeContains(&Contains)
	MakeContains(&StringContains)
	// MakeMap(&Map)
	MakeMap(&StringMap)
	MakeMap(&StringToBoolMap)
	MakePartition(&PartitionInt)
}

var Contains func(interface{}, interface{}) bool

var StringContains func([]string, string) bool

// var Map func(interface{}, func(interface{}) interface{}) interface{}

var StringMap func([]string, func(string) string) []string

var StringToBoolMap func([]string, func(string) bool) []bool

var Partition func(interface{}, func(interface{}) bool) interface{}

var PartitionInt func([]int, func(int) bool) ([]int, []int)

var PartitionString func([]string, func(string) bool) []string


func MakeX(wrapper interface{}, fn func(args []reflect.Value) (results []reflect.Value)) {
	wrapperFn := reflect.ValueOf(wrapper).Elem()
	v := reflect.MakeFunc(wrapperFn.Type(), fn)
	wrapperFn.Set(v)
}

func MakeContains(fn interface{}) {
	MakeX(fn, _contains)
}

func MakeMap(fn interface{}) {
	MakeX(fn, _map)
}

func MakePartition(fn interface{}) {
	MakeX(fn, _partition)
}

func _contains(values []reflect.Value) []reflect.Value {

	v := interfaceToValue(values[0])
	o := values[1].Interface()

	for i := 0; i < v.Len(); i++ {
		e := v.Index(i).Interface()
		if e == o {
			return wrap(reflect.ValueOf(true))
		}
	}
	return wrap(reflect.ValueOf(false))
}

func _map(values []reflect.Value) []reflect.Value {

	v := interfaceToValue(values[0])
	fn := values[1]

	outType := reflect.SliceOf(fn.Type().Out(0))
	ret := reflect.MakeSlice(outType, v.Len(), v.Len())

	for i := 0; i < v.Len(); i++ {
		e := v.Index(i)
		r := fn.Call([]reflect.Value{e})
		ret.Index(i).Set(r[0])
	}
	return wrap(ret)
}


func _partition(values []reflect.Value) []reflect.Value {
	slice := values[0]
	fn := values[1]

	t := reflect.MakeSlice(slice.Type(), 0, 0)
	f := reflect.MakeSlice(slice.Type(), 0, 0)

	for i := 0; i < slice.Len(); i++ {
		e := slice.Index(i)
		r := fn.Call([]reflect.Value{e})
		if r[0].Bool() {
			t = reflect.Append(t, e)
		} else {
			f = reflect.Append(f, e)
		}
	}
	return []reflect.Value{t, f}
}

func wrap(v reflect.Value) []reflect.Value {
	return []reflect.Value{v}
}

func interfaceToValue(v reflect.Value) reflect.Value {
	if v.Kind() == reflect.Interface {
		return reflect.ValueOf(v.Interface())
	}
	return v
}





func partition (slice []int, fn func(int) bool) ([]int, []int) {
	a := []int{}
	b := []int{}

	for i := 0; i < len(slice); i++ {
		e := slice[i]
		if fn(e) {
			a = append(a, e)
		} else {
			b = append(b, e)
		}
	}

	return a, b
}

