package kafka

import (
	"testing"
	"context"
	config2 "github.com/modeyang/LogCruiser/config"
	"time"
	"log"
)

func TestKafkaSink(T *testing.T){
	config := config2.ConfigRaw {
		"type": "kafka",
		"topic": "test_go",
		"brokers": "10.100.4.149:9092",
	}
	ctx , cancel := context.WithCancel(context.Background())
	defer cancel()

	sinker, err := InitSinkHandler(ctx, &config)
	if err != nil {
		T.Error(err)
		return
	}
	timestamp := time.Now().Unix()
	timestamp = timestamp - timestamp % 60
	metric := config2.MetricResult{
		Timestamp: timestamp,
		Data: map[string]int64{
			"error.qps/host=All": 1,
		},
	}
	log.Println(metric)
	err = sinker.Push(ctx, metric)
	if err != nil {
		T.Error(err)
		return
	}
}
