package convert

import (
	"testing"
	"github.com/modeyang/LogCruiser/config"
	"log"
	"github.com/modeyang/LogCruiser/config/logevent"
	"time"
	"context"
	"reflect"
)

func TestConvertFilterHandler(t *testing.T) {
	rawConfig := config.ConfigRaw{
		"fields": []config.ConfigRaw{
			{
				"field": "status",
				"to": "int",
				"removeIfFail": true,
				"multiplier": 10,
			},
			{
				"field": "request_time",
				"to": "float",
				"removeIfFail": true,
				"multiplier": 1000,
			},
		},
	}

	conf := defaultFilterConfig()
	if err := config.ReflectConfig(&rawConfig, &conf); err != nil {
		t.Error(err)
		t.Error("reflect config failed")
	}
	log.Println(conf)

	event := logevent.LogEvent{
		Timestamp: time.Now(),
		Tags: []string{},
		Event:map[string]interface{}{
			"status": "500",
			"request_time": "0.5",
		},
	}
	ctx := context.Background()
	event = conf.Event(ctx, event)
	log.Println(event)
	if event.Event["status"] != int64(5000) {
		t.Error("error convert")
	}
	log.Println(reflect.TypeOf(event.Event["request_time"]))
	if event.Event["request_time"] != int64(500) {
		t.Error("error convert request_time failed")
	}
}
