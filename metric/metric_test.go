package metric

import (
	"testing"
	"text/template"
	"bytes"
	"log"
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

