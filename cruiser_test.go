package main

import (
	"testing"
	"reflect"
)

func Test_reflect(t *testing.T) {
	type person struct {
		name string
	}
	p := person{name:"haha"}
	rt := reflect.TypeOf(p)
	if rt.Kind() != reflect.Struct {
		t.Error(p, " is not struct type")
	}
}
