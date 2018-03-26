package metric

import (
	"github.com/rcrowley/go-metrics"
	"gopkg.in/yaml.v2"
	"bytes"
	"text/template"
	"strconv"
	"log"
	"sync"
	"time"
	"github.com/modeyang/LogCruiser/config/logevent"
	"github.com/modeyang/LogCruiser/config"
)

type MetricItemConfig struct {
	MetricTmpl string `yaml:"metric"`
	MetricType string `yaml:"type"`
	MetricValue interface{} `yaml:"value"`
	FilterTmpls []string `yaml:"filters"`

	MetricName 		*template.Template
	FilterFuncs	   	[]*template.Template
}

var NAMESPACE = "Metric"

func NewMetricFromConfig(data []byte, metric_item *MetricItemConfig)(error) {
	err := yaml.Unmarshal(data, metric_item)
	if err != nil {
		return err
	}
	tmpl, err := template.New(NAMESPACE).Parse(metric_item.MetricTmpl)
	if err != nil {
		return err
	}
	metric_item.MetricName = tmpl
	switch metric_item.MetricValue.(type) {
	case string:
		metric_item.MetricValue, _ = template.New("test").Parse(metric_item.MetricValue.(string))
	}
	if len(metric_item.FilterTmpls) > 0 {
		metric_item.FilterFuncs = []*template.Template{}
		for i, tpl := range(metric_item.FilterTmpls) {
			tmpl, err := template.New(string(i)).Parse(tpl)
			if err != nil {
				log.Fatalln(err)
				return err
			}
			metric_item.FilterFuncs = append(metric_item.FilterFuncs, tmpl)
		}
	}
	return nil
}

func RenderTemplate(tmpl *template.Template, event map[string]interface{}) (string, error) {
	var tpl bytes.Buffer
	err := tmpl.Execute(&tpl, event)
	if err != nil {
		log.Fatalln(err)
		return "", err
	}
	return tpl.String(), nil
}

func (mtr *MetricItemConfig)Render(event map[string]interface{})(string, error) {
	return RenderTemplate(mtr.MetricName, event)
}

func (mtr *MetricItemConfig)RenderValue(event map[string]interface{})int64{
	switch mtr.MetricValue.(type) {
	case int:
		return int64(mtr.MetricValue.(int))
	case *template.Template:
		value ,_:= RenderTemplate(mtr.MetricValue.(*template.Template), event)
		if value != "" {
			intValue, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				panic(err)
			}
			return intValue
		}
	}
	panic(mtr.MetricValue)
}

func (mtr *MetricItemConfig)Filter(event map[string]interface{})bool {
	for _, filter_tmpl := range(mtr.FilterFuncs) {
		pass, _ := RenderTemplate(filter_tmpl, event)
		if pass == "" {
			return false
		}
	}
	return true
}

type MetricResult struct {
	Timestamp time.Time
	AppMetric map[string]int64
}

type MetricConfig struct {
	Metrics 	[]*MetricItemConfig 		`yaml:"metrics"`
	MapResults 	map[int64]MetricResult

	MapRegistry map[int64]metrics.StandardRegistry
}

type MetricChan chan MetricResult

//func NewMetricProcessor(m []*MetricItemConfig)*MetricConfig{
//	return &MetricConfig{Metrics: m, MapResults: make(map[int64]MetricResult)}
//}

func (c *config.Config) startCalcuateMetrics() {
	c.eg.Go(func()error{
		return nil
	})

}

var MetricRegistry = metrics.DefaultRegistry
//
//func (m *MetricConfig)Calculate(logEvent *logevent.LogEvent)error{
//	var wg sync.WaitGroup
//	event := logEvent.Event
//	for _, mtr := range(m.Metrics) {
//		metric_name, err := mtr.Render(event)
//		if err != nil {
//			log.Printf("render template %v failed\n", mtr.MetricTmpl)
//			continue
//		}
//		wg.Add(1)
//		go func(this *MetricConfig, mtr *MetricItemConfig, event map[string]interface{}){
//			if mtr.Filter(event) {
//				switch {
//				case mtr.MetricType == "counter" || mtr.MetricType == "c":
//					metricFunc := metrics.GetOrRegisterCounter(metric_name, MetricRegistry)
//					metricFunc.Inc(mtr.RenderValue(event))
//				default:
//					log.Printf("metric type %s is not right\n", mtr.MetricType)
//				}
//			}
//			wg.Done()
//		}(m, mtr, event)
//
//	}
//	wg.Wait()
//	return nil
//}
//
//
//
//func (m *MetricConfig)GetMetrics()interface{}{
//	raw_metrics := map[string]int64{}
//	allMetrics := MetricRegistry.GetAll()
//	for k, v := range(allMetrics) {
//		if c, ok := v["count"]; ok {
//			raw_metrics[k] = c.(int64)
//		}
//	}
//	MetricRegistry.UnregisterAll()
//	return raw_metrics
//}

