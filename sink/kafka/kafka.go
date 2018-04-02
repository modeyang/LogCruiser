package kafka

import (
	"github.com/modeyang/LogCruiser/config"
	"context"
	"log"
	"github.com/Shopify/sarama"
	"strings"
	"fmt"
	"encoding/json"
)


type KafkaSinkConfig struct {
	config.CommonConfig

	Topic 			string 		`yaml:"topic"`
	Brokers 		string 		`yaml:"brokers"`


	Producer 	sarama.SyncProducer
}

const Module = "kafka"

func defaultSinkConfig() KafkaSinkConfig {
	return KafkaSinkConfig{
		CommonConfig: config.CommonConfig{
			Type: Module,
		},
	}
}

func(sink *KafkaSinkConfig)Push(ctx context.Context, result config.MetricResult)error{
	jsonMetric, err := json.Marshal(result)
	if err != nil {
		return err
	}
	msg := &sarama.ProducerMessage{
		Topic: sink.Topic,
		Value: sarama.StringEncoder(jsonMetric),
	}
	partition, offset, err := sink.Producer.SendMessage(msg)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println(fmt.Sprintf("send metric success, topic: %s partition: %v offset: %v", sink.Topic, partition, offset))
	return nil
}


func InitSinkHandler(ctx context.Context, raw *config.ConfigRaw)(config.TypeSinkConfig, error) {
	conf := defaultSinkConfig()
	if err := config.ReflectConfig(raw, &conf); err != nil {
		return nil, err
	}

	producer, err := sarama.NewSyncProducer(strings.Split(conf.Brokers, ","), nil)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	conf.Producer = producer
	return &conf, nil
}

