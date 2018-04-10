package config

import (
	"context"
	"github.com/modeyang/LogCruiser/config/logevent"
	"errors"
)

type TypeInputConfig interface {
	TypeCommonConfig
	Source(ctx context.Context, msgChan chan <- logevent.LogEvent)(err error)
}

type InputConfig struct {
	CommonConfig
}

type InputHandler func(ctx context.Context, raw *ConfigRaw)(TypeInputConfig, error)

var (
	mapInputHandler = map[string]InputHandler{}
)

func RegisterInputHandler(name string, handler InputHandler){
	mapInputHandler[name] = handler
}

func (c *Config) getInputs()(inputs []TypeInputConfig, err error) {
	var input TypeInputConfig
	for _, raw := range(c.InputRaw) {
		handler, ok := mapInputHandler[raw["type"].(string)]
		if ! ok {
			return inputs, errors.New("unknown input type: " + raw["type"].(string))
		}
		if input, err = handler(c.ctx, &raw); err != nil {
			return inputs, errors.New("init input module failed : " + raw["type"].(string))
		}
		inputs = append(inputs, input)
	}
	return
}

func (c *Config) startInput()(error) {
	inputs, err := c.getInputs()
	if err != nil {
		return err
	}
	for _, input := range(inputs) {
		go func(input TypeInputConfig) {
			c.eg.Go(func() error {
				return input.Source(c.ctx, c.chInFilter)
			})
		}(input)
	}
	return nil
}