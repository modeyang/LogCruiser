package kafka

import (
	"testing"
	"context"
	"time"
	config2 "github.com/modeyang/LogCruiser/config"
	"fmt"
)

func TestKafkaSource(T *testing.T) {
	config := config2.ConfigRaw {
		"type": "kafka",
		"topics": []string{ "ops-https-accesslog" },
		"brokers": "kafka-10-100-4-129:9092,kafka-10-100-4-135:9092,kafka-10-100-4-136:9092,kafka-10-100-4-137:9092",
		"group": "ops_https-slog-go",
	}

	ctx , cancel := context.WithCancel(context.Background())
	input, err := KafkaInputHandler(ctx, &config)
	if err != nil {
		T.Error(err)
	}
	msgChan := make(config2.MsgChan, 10)
	timer := time.NewTimer(3)
	go func() {
		for _ = range(timer.C) {
			cancel()
		}
	}()

	input.Source(ctx, msgChan)

	fmt.Println("input end")
	for msg := range(msgChan) {
		fmt.Println(msg)
	}

}