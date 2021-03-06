package kafka

import (
	"testing"
	"context"
	"time"
	config2 "github.com/modeyang/LogCruiser/config"
	"log"
)

func TestKafkaSource(T *testing.T) {
	config := config2.ConfigRaw {
		"type": "kafka",
		"topics": []string{ "ops-https-accesslog" },
		"brokers": "10.100.4.149:9092",
		"group": "ops_https-slog-go",
	}

	ctx , cancel := context.WithCancel(context.Background())
	defer cancel()
	input, err := KafkaInputHandler(ctx, &config)
	if err != nil {
		T.Error(err)
	}
	msgChan := make(config2.MsgChan, 10)
	timer := time.NewTimer(1 * time.Second)
	log.Println("timer start")

	go input.Source(ctx, msgChan)

	log.Println("input end")
	for {
		select {
		case msg, ok := <- msgChan:
			if ok {
				log.Println(msg)
			}
		case <- timer.C:
			log.Println("timer stop")
			goto end
		}
	}
end:
	log.Println("print msg end")

}