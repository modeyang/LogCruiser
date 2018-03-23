package main

import (
	"encoding/json"
	"github.com/modeyang/LogCruiser/filter"
	"github.com/modeyang/LogCruiser/metric"
	"log"
	"time"
	"fmt"
	"runtime"
)


var LogFilters []filter.LogFilter
func init() {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)

	LogFilters = make([]filter.LogFilter, 0, 5)
	fields := []string{"fromhost", "idc", "timestamp", "count"}
	remove_fields := []string{"message"}
	split_filter := filter.NewSplitFilter("message", fields, "|", remove_fields)
	LogFilters = append(LogFilters, split_filter)
}

var CONFIG = `
metric: "access.qps/fromhost={{.fromhost}}"
type: "c"
value: "{{ .count }}"
`

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	start := time.Now()
	event := make(map[string]interface{})
	var msg = "10.100.1.145|bjyg|21/Mar/2018:17:24:12 +0800|2"
	event["message"] = msg

	for _, filterFun := range(LogFilters) {
		event, _ = filterFun.Filter(event)
	}
	log.Println(event)
	var metric_item metric.MetricItem
	err := metric.NewMetricFromConfig([]byte(CONFIG), &metric_item)
	if err != nil {
		log.Println(err)
		return
	}
	metricResults := metric.NewMetricResult([]*metric.MetricItem{&metric_item},)
	var jm []byte
	index := 0
	timer := time.NewTimer(3 * time.Second )
	go func() {
		for i:=0; i < 100000; i++ {
			metricResults.Calculate(event)
			jm, _ = json.Marshal(metricResults.GetMetrics())
			index ++
			//log.Println(string(jm))
		}

	}()
	select {
	case <- timer.C:
		log.Println("timeout")
	}
	log.Println(index)
	log.Println(string(jm))
	t := time.Now()
	elapsed := t.Sub(start)
	fmt.Println(elapsed)
}
