package main

import (
	"fmt"
	"text/template"
	"bytes"
	"encoding/json"
	"github.com/modeyang/LogCruiser/filter"
)
//type MetricFunc struct {
//	MetricName string
//}
//
//type MetricProcessor interface {
//	Calcuate(event map[string]interface{})(interface{}, error)
//	Filter(event map[string]interface{})bool
//}
//
//type Metrics struct {
//	Result map[string]interface{}
//	MtrFuncs []MetricFunc
//}
//
//func (mi *MetricFunc) Calcuate(event map[string]interface{}) {
//
//}
//
//func (mi *MetricFunc) Filter(event map[string]interface{}) bool{
//	return true
//}
//
//func (mtr *Metrics) DoCalcuate(event map[string]interface{}) {
//
//}

func render(tmpl *template.Template, event map[string]interface{}) (string, error) {
	var tpl bytes.Buffer
	err := tmpl.Execute(&tpl, event)
	if err != nil {
		return "", err
	}
	return tpl.String(), nil
}

var LogFilters []filter.LogFilter
func init() {
	fmt.Println("do in init1")
	LogFilters = make([]filter.LogFilter, 0, 5)
	fields := []string{"fromhost", "syslogt	ags", "timestamp"}
	remove_fields := []string{"message"}
	split_filter := filter.NewSplitFilter("message", fields, "|", remove_fields)
	LogFilters = append(LogFilters, split_filter)
}

func main() {
	event := make(map[string]interface{})
	var msg = "10.100.1.145|bjyg,odin|21/Mar/2018:17:24:12 +0800"
	event["message"] = msg

	for _, filter := range(LogFilters) {
		event, _ = filter.Filter(event)
	}
	fmt.Println(event)

	results := make(map[string]interface{})
	tmpl, err := template.New("access").Parse("access.qps/fromhost={{.fromhost}}")
	if err != nil {
		fmt.Println(err)
		return
	}
	metric_name, err := render(tmpl, event)
	if err != nil {
		fmt.Println(err)
		return
	}
	if _, ok := results[metric_name]; ! ok {
		results[metric_name] = 1
	}else {
		results[metric_name] = results[metric_name].(int)
	}
	fmt.Print(results)
	jresults, _ := json.Marshal(results)
	fmt.Println(string(jresults))
}
