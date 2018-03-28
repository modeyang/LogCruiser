package config

import (
	"context"
	"github.com/modeyang/LogCruiser/config/logevent"
	"github.com/modeyang/LogCruiser/metric"
	"errors"
	"github.com/rcrowley/go-metrics"
)

type TypeMetricConfig interface {
	TypeCommonConfig
	Calculate(ctx context.Context, registry metrics.Registry, event logevent.LogEvent)error
}

type MetricResult struct {
	Timestamp int 				`json:"timestamp"`
	Data 	  map[string]int64 	`json:"data"`
}


func (c *Config)getMetrics()(metrics []TypeMetricConfig, err error) {
	var metricItem TypeMetricConfig
	for _, raw := range(c.MetricRaw){
		metricItem, err = metric.InitMetricConfig(c.ctx, &raw)
		if err != nil {
			return metrics, errors.New("unable init metric : " + Map2String(raw))
		}
		metrics = append(metrics, metricItem)
	}
	return
}

func (c *Config)startMetrics()(err error){
	allMetrics, err := c.getMetrics()
	if err != nil {
		return err
	}
	for {
		select {
		case <- c.ctx.Done():
			return nil
		case event := <- c.chInMetric:
			for _, metricItem := range(allMetrics) {
				func(item TypeMetricConfig)error{
					c.eg.Go(func() error {
						return item.Calculate(c.ctx, c.registry, event)
					})
					return nil
				}(metricItem)

			}
		}
	}
	return nil
}