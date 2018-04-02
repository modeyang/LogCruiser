package kafka

import (
	"github.com/modeyang/LogCruiser/config"
	"context"
	"github.com/modeyang/LogCruiser/config/logevent"
	"github.com/Shopify/sarama"
	cluster "github.com/bsm/sarama-cluster"
	"errors"
	"strings"
	"log"
	"fmt"
)

type InputConfig struct {
	config.InputConfig

	Topics 		[]string  	`yaml:"topics"`
	Brokers 	string		`yaml:"brokers"`
	Group		string 		`yaml:"group"`
	Offset 		string 		`yaml:"offset"`
	ReturnError bool		`yaml:"return_error"`
	ReturnNotify bool 		`yaml:"return_notify"`

	Consumer   *cluster.Consumer
}

const Module = "kafka"

func defaultInputConfig()InputConfig{
	return InputConfig{
		InputConfig: config.InputConfig{
			CommonConfig: config.CommonConfig{
				Type:Module,
			},
		},
		ReturnError: true,
		ReturnNotify: true,
		Offset: "newest",
	}
}

func KafkaInputHandler(ctx context.Context, raw *config.ConfigRaw)(config.TypeInputConfig, error) {
	conf := defaultInputConfig()
	err := config.ReflectConfig(raw, &conf)
	if err != nil {
		return nil, err
	}
	if conf.Offset != "newest" && conf.Offset != "oldest" {
		return nil, errors.New(fmt.Sprintf("offset: %s not in [newest, oldest]", conf.Offset))
	}
	kcfg := cluster.NewConfig()
	kcfg.Consumer.Return.Errors = true
	kcfg.Group.Return.Notifications = true
	if conf.Offset == "newest" {
		kcfg.Consumer.Offsets.Initial = sarama.OffsetNewest
	} else {
		kcfg.Consumer.Offsets.Initial = sarama.OffsetOldest
	}
	consumer, err := cluster.NewConsumer(strings.Split(conf.Brokers, ","), conf.Group, conf.Topics, kcfg)
	if err != nil {
		return nil, err
	}
	conf.Consumer = consumer
	return &conf, nil
}

func (input *InputConfig)Source(ctx context.Context, msgChan chan <- logevent.LogEvent)(err error){
	for {
		select {
		case <- ctx.Done():
			return nil
		case err, ok := <- input.Consumer.Errors():
			if ok  {
				log.Printf("Error: %s\n", err.Error())
			}
		case ntf, ok := <- input.Consumer.Notifications():
			if ok {
				log.Printf("Rebalanced: %+v\n", ntf)
			}
		case msg := <- input.Consumer.Messages():
			msgChan <- *logevent.NewLogEvent(string(msg.Value))
			input.Consumer.MarkOffset(msg, "")
		}
	}
	return nil
}
