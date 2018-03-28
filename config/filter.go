package config

import (
	"context"
	"github.com/modeyang/LogCruiser/config/logevent"
	"errors"
)

type TypeFilterConfig interface {
	TypeCommonConfig
	Event(context.Context, logevent.LogEvent)logevent.LogEvent
}

type FilterConfig struct {
	CommonConfig
}

// filter handler
type FilterHandler func(ctx context.Context, raw *ConfigRaw)(TypeFilterConfig, error)

var mapFilterHandler = map[string]FilterHandler{}


func RegisterFilterHandler(name string, handler FilterHandler) {
	mapFilterHandler[name] = handler
}

func (c *Config) getFilters()(filters []TypeFilterConfig, err error) {
	var filter TypeFilterConfig
	for _, raw := range(c.FilterRaw) {
		handler, ok := mapFilterHandler[raw["type"].(string)]
		if ! ok {
			return filters, errors.New("unknown filter type: " + raw["type"].(string))
		}
		if filter, err = handler(c.ctx, &raw); err != nil {
			return filters, errors.New("init filter module failed : " + raw["type"].(string))
		}
		filters = append(filters, filter)
	}
	return
}

func (c *Config)startFilters() (err error) {
	filters, err := c.getFilters()
	if err != nil {
		return
	}
	c.eg.Go(func() error {
		for {
			select {
			case <- c.ctx.Done():
				return nil
			case event := <- c.chInFilter:
				for _, filter := range(filters) {
					event = filter.Event(c.ctx, event)
				}
				c.chInMetric <- event
			}
		}
	})
	return
}

