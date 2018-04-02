package module

import (
	"github.com/modeyang/LogCruiser/config"
	"github.com/modeyang/LogCruiser/input/file"
	"github.com/modeyang/LogCruiser/filter/split"
	"github.com/modeyang/LogCruiser/filter/convert"
	"github.com/modeyang/LogCruiser/sink/debugsink"
	"github.com/modeyang/LogCruiser/metric"
	"github.com/modeyang/LogCruiser/input/kafka"
	sinkKafka "github.com/modeyang/LogCruiser/sink/kafka"
)

func InitModule(){
	config.RegisterInputHandler(file.Module, file.FileInputHandler)
	config.RegisterInputHandler(kafka.Module, kafka.KafkaInputHandler)

	config.RegisterFilterHandler(split.Module, split.SplitFilterHandler)
	config.RegisterFilterHandler(convert.Module, convert.ConvertFilterHandler)

	config.RegisterSinkHandler(debugsink.Module, debugsink.InitSinkHandler)
	config.RegisterSinkHandler(sinkKafka.Module, sinkKafka.InitSinkHandler)

	config.RegisterMetricHandlers(metric.InitMetricConfig)
}
