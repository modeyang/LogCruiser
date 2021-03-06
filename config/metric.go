package config

import (
	"context"
	"github.com/modeyang/LogCruiser/config/logevent"
	"errors"
	"github.com/rcrowley/go-metrics"
	"log"
)

type TypeMetricConfig interface {
	TypeCommonConfig
	Calculate(ctx context.Context, registry metrics.Registry, event logevent.LogEvent)error
}

type MetricResult struct {
	Timestamp int64 			`json:"timestamp"`
	Data 	  map[string]int64 	`json:"data"`
}

var LogProcessMetric = "log_process.count"
var listMetricHanlder = []HandlerMetric{}

type HandlerMetric func(ctx context.Context, raws *ConfigRaw)(TypeMetricConfig, error)

func RegisterMetricHandlers(handler HandlerMetric) {
	listMetricHanlder = append(listMetricHanlder, handler)
}

func (c *Config)getMetrics()(metrics []TypeMetricConfig, err error) {
	var metricItem TypeMetricConfig
	if len(listMetricHanlder) == 0 {
		return metrics, errors.New("no metric handler")
	}
	handler := listMetricHanlder[0]
	for _, raw := range(c.MetricRaw){
		metricItem, err = handler(c.ctx, &raw)
		if err != nil {
			return metrics, errors.New("unable init metric : " + Map2String(raw))
		}
		metrics = append(metrics, metricItem)
	}
	return
}

func (c *Config)handleSelfMetric(name string){
	metrics.GetOrRegisterCounter(name, c.selfRegistry).Inc(1)
}

func (c *Config)startMetrics()(err error){
	log.Println("start metrics")
	allMetrics, err := c.getMetrics()
	if err != nil {
		return err
	}
	c.eg.Go(func() error {
		for {
			select {
			case <- c.ctx.Done():
				return nil
			case event := <- c.chInMetric:
				c.handleSelfMetric(LogProcessMetric)
				for _, metricItem := range(allMetrics) {
					metricItem.Calculate(c.ctx, c.registry, event)
				}
			}
		}
	})
	return nil
}