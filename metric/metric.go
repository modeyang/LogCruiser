package metric

import (
	"text/template"
	"strconv"
	"log"
	"github.com/modeyang/LogCruiser/config"
	"context"
	"github.com/modeyang/LogCruiser/config/logevent"
	"github.com/rcrowley/go-metrics"
	"reflect"
	"strings"
)

const ERROR_STRING = "<no value>"

type MetricConfig struct {
	config.CommonConfig
	MetricTmpl 	string 		`yaml:"metricTmpl"`
	MetricValue interface{} `yaml:"metricValue"`
	FilterTmpls []string 	`yaml:"filterTmpls,omitempty"`

	MetricName 		*template.Template
	FilterFuncs	   	[]*template.Template
}

var NAMESPACE = "Metric"

func DefaultMetricConfig() MetricConfig {
	return MetricConfig{
		CommonConfig: config.CommonConfig{
			Type: "counter",
		},
	}
}

func InitMetricConfig(ctx context.Context, raw *config.ConfigRaw) (config.TypeMetricConfig, error) {
	conf := DefaultMetricConfig()
	err := config.ReflectConfig(raw, &conf)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	tmpl, err := template.New(NAMESPACE).Parse(conf.MetricTmpl)
	if err != nil {
		return nil, err
	}
	conf.MetricName = tmpl
	switch conf.MetricValue.(type) {
	case string:
		conf.MetricValue, _ = template.New("test").Parse(conf.MetricValue.(string))
	}
	if len(conf.FilterTmpls) > 0 {
		conf.FilterFuncs = []*template.Template{}
		for i, tpl := range(conf.FilterTmpls) {
			tmpl, err := template.New(string(i)).Parse(tpl)
			if err != nil {
				log.Println(err)
				return nil, err
			}
			conf.FilterFuncs = append(conf.FilterFuncs, tmpl)
		}
	}
	return &conf, nil
}


func (mtr *MetricConfig)render(event map[string]interface{})(string, error) {
	return config.RenderTemplate(mtr.MetricName, event)
}

func (mtr *MetricConfig)renderValue(event map[string]interface{})int64{
	switch mtr.MetricValue.(type) {
	case int:
		return int64(mtr.MetricValue.(int))
	case float64:
		return int64(mtr.MetricValue.(float64))
	case *template.Template:
		value ,_:= config.RenderTemplate(mtr.MetricValue.(*template.Template), event)
		if value != ERROR_STRING && len(value) > 0 {
			//if strings.Contains(value, ".") {
			//	floatValue, err:= strconv.ParseFloat(value, 64)
			//	if err != nil {
			//		log.Println(err)
			//		return 0
			//	}
			//	return floatValue
			//}
			if strings.Contains(value, ".") {
				floatValue, err := strconv.ParseFloat(value, 64)
				if err != nil {
					log.Println(err)
					return 0
				}
				return int64(floatValue)
			}
			intValue, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				log.Println(err)
				return 0
			}
			return intValue
		}
	default:
		log.Println(reflect.TypeOf(mtr.MetricValue))
	}
	return 0
}

func (mtr *MetricConfig)filter(event map[string]interface{})bool {
	for _, filter_tmpl := range(mtr.FilterFuncs) {
		pass, _ := config.RenderTemplate(filter_tmpl, event)
		if pass == "" {
			return false
		}
	}
	return true
}

func (mtr *MetricConfig)Calculate(ctx context.Context, registry metrics.Registry, event logevent.LogEvent)error {
	metric_name, err := mtr.render(event.Event)
	if err != nil || metric_name == ERROR_STRING {
		log.Printf("render template %v failed\n", mtr.MetricTmpl)
		return nil
	}
	if ok := mtr.filter(event.Event); ok {
		switch  {
		case mtr.Type == "counter" || mtr.Type == "c":
			metricFunc := metrics.GetOrRegisterCounter(metric_name, registry)
			metricFunc.Inc(mtr.renderValue(event.Event))
		default:
			log.Printf("metric type %s is not right\n", mtr.Type)
		}
	}
	return nil
}


