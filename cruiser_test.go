package main

import (
	"testing"
	"reflect"
	"gopkg.in/yaml.v2"
	"log"
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

func Test_Yaml(t *testing.T) {
	type TestYaml struct {
		Name 		string `yaml:"name"`
		FirstName 	string `yaml:"first_name"`
	}

data := `
name: hehe
first_name: fda
`
	result := TestYaml{}
	err := yaml.Unmarshal([]byte(data), &result)
	if err != nil {
		t.Error(err)
	}
	log.Println(result.FirstName)
}
