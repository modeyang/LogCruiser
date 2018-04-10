package debugsink

import (
	"testing"
	"time"
	"context"
	config2 "github.com/modeyang/LogCruiser/config"
)

func TestDebug(T *testing.T) {
	config := config2.ConfigRaw {
		"type": "kafka",
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
	err = sinker.Push(ctx, metric)
	if err != nil {
		T.Error(err)
		return
	}
}
