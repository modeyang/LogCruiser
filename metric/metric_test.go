package metric

import (
	"testing"
	"text/template"
	"bytes"
	"log"
	"fmt"
)


func TestMetricResult_GetMetrics(t *testing.T) {
	s := "{{ if ge .status 500 }}true{{ end }}"
	tmpl, err:= template.New("test").Parse(s)
	if err != nil {
		t.Error("can not template string " + s)
	}
	event := map[string]interface{}{"status": 500}
	var tpl bytes.Buffer
	err = tmpl.Execute(&tpl, event)
	if err != nil {
		t.Error("render template failed: " + tmpl.Name())
		log.Println(err)
	}
	if tpl.String() != "true" {
		t.Error(tpl.String() + "not equal true")
	}
}

func TestTemplate(t *testing.T) {
	s := "{{ .status }}"
	tmpl, err:= template.New("test").Parse(s)
	if err != nil {
		t.Error("can not template string " + s)
	}
	event := map[string]interface{}{"status": 500}
	var tpl bytes.Buffer
	err = tmpl.Execute(&tpl, event)
	if err != nil {
		t.Error("render template failed: " + tmpl.Name())
		log.Println(err)
	}
	if tpl.String() != "500" {
		t.Error(tpl.String() + " not equal 500 + \n")
	}

	var b bytes.Buffer
	fmt.Println(b.Len())

	//no status field
	event1 := map[string]interface{}{"haha": 200}
	err = tmpl.Execute(&b, event1)
	if err != nil {
		t.Error("render template failed: " + tmpl.Name())
		log.Println(err)
	}
	fmt.Println(b.Len())
	fmt.Println(b.String())
	//if b.Len() != 0 {
	//	t.Error(b.String() + " not equal 200")
	//}

	if b.String() != "<no value>" {
		t.Error(b.String() + " not equal 200")
	}

}

