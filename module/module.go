package module

import (
	"github.com/modeyang/LogCruiser/config"
	"github.com/modeyang/LogCruiser/input/file"
	"github.com/modeyang/LogCruiser/filter/split"
	"github.com/modeyang/LogCruiser/filter/convert"
	"github.com/modeyang/LogCruiser/sink/debugsink"
	"github.com/modeyang/LogCruiser/metric"
)

func InitModule(){
	config.RegisterInputHandler(file.Module, file.FileInputHandler)

	config.RegisterFilterHandler(split.Module, split.SplitFilterHandler)
	config.RegisterFilterHandler(convert.Module, convert.ConvertFilterHandler)

	config.RegisterSinkHandler(debugsink.Module, debugsink.InitSinkHandler)

	config.RegisterMetricHandlers(metric.InitMetricConfig)
}
