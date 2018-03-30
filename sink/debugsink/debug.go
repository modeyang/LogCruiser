package debugsink

import (
	"github.com/modeyang/LogCruiser/config"
	"context"
	"log"
)

type DebugSinkConfig struct {
	config.CommonConfig
}

const Module = "debug"

func(d *DebugSinkConfig)Push(ctx context.Context, result config.MetricResult)error{
	log.Println(result)
	return nil
}


func InitSinkHandler(ctx context.Context, raw *config.ConfigRaw)(config.TypeSinkConfig, error) {
	conf := DebugSinkConfig{
		CommonConfig: config.CommonConfig{Type:"debug"},
	}
	if err := config.ReflectConfig(raw, &conf); err != nil {
		return nil, err
	}
	return &conf, nil
}

