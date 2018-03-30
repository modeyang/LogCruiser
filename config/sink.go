package config

import (
	"context"
	"errors"
	"log"
)

type TypeSinkConfig interface {
	TypeCommonConfig

	Push(context.Context, MetricResult)error
}


type SinkHandler func(context.Context, *ConfigRaw)(TypeSinkConfig, error)

var mapSinkHandler = map[string]SinkHandler{}

func RegisterSinkHandler(name string, handler SinkHandler) {
	mapSinkHandler[name] = handler
}


func (c *Config)getSinkers()(sinkers []TypeSinkConfig, err error) {
	var sinker TypeSinkConfig
	for _, raw := range(c.SinkRaw) {
		handler, ok := mapSinkHandler[raw["type"].(string)]
		if !ok {
			return sinkers, errors.New("unknown sinker type " + raw["type"].(string))
		}
		if sinker, err = handler(c.ctx, &raw); err != nil {
			return sinkers, errors.New("init sinker module failed : " + raw["type"].(string))
		}
		sinkers = append(sinkers, sinker)
	}
	return
}

func (c *Config)startSinkers()error {
	log.Println("start sinkers")
	sinkers, err := c.getSinkers()
	if err != nil {
		return err
	}
	for {
		select {
		case <- c.ctx.Done():
			return nil
		case metricResult := <- c.chInSinker:
			for _, sink := range(sinkers) {
				func(s TypeSinkConfig)error{
					c.eg.Go(func() error {
						return s.Push(c.ctx, metricResult)
					})
					return nil
				}(sink)
			}
		}
	}
}
