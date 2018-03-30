package split

import (
	"strings"
	"github.com/modeyang/LogCruiser/config"
	"github.com/modeyang/LogCruiser/config/logevent"
	"context"
)

type FilterConfig struct {
	config.CommonConfig

	Source	string  `yaml:"source"`
	Fields []string `yaml:"fields"`
	SplitField 	string 	`yaml:"splitField"`
	RemoveFields []string `yaml:"removeFields:omitempty"`
}

const Module = "split"

func defaultFilterConfig()FilterConfig{
	return FilterConfig{
		CommonConfig: config.CommonConfig{Type:Module},
		Source: "message",
		Fields: []string{},
		RemoveFields: []string{},
	}
}

func SplitFilterHandler(ctx context.Context, raw *config.ConfigRaw)(config.TypeFilterConfig, error){
	conf := defaultFilterConfig()
	if err := config.ReflectConfig(raw, &conf); err != nil {
		return nil, err
	}
	return &conf, nil
}

func (f *FilterConfig)Event(ctx context.Context,event logevent.LogEvent)logevent.LogEvent{
	if value, ok := event.Event[f.Source]; ok {
		msg_list := strings.Split(value.(string), f.SplitField)
		for i, f := range(f.Fields) {
			if i <= len(msg_list) - 1 {
				event.Event[f] = string(msg_list[i])
			} else {
				event.Event[f] = ""
			}
		}
		if len(f.RemoveFields) > 0 {
			for _, k := range(f.RemoveFields) {
				delete(event.Event, k)
			}
		}
	}
	return event
}

