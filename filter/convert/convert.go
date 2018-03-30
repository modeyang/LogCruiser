package convert

import (
	"github.com/modeyang/LogCruiser/config"
	"github.com/modeyang/LogCruiser/config/logevent"
	"context"
	"strconv"
	"log"
)

type convertItem struct {
	Field 			string 	`yaml:"field"`
	To 				string 	`yaml:"to"`
	RemoveIfFail 	bool 	`yaml:"removeIfFail,omitempty"`
	Multiplier		int 	`yaml:"multiplier,omitempty"`
}

type T interface{}

func NewConvertItem(raw config.ConfigRaw) (convertItem, error){
	item := convertItem{}
	err := config.ReflectConfig(&raw, &item)
	return item, err
}

type FilterConfig struct {
	config.CommonConfig

	Fields []config.ConfigRaw `yaml:"fields"`
}

const Module = "convert"

type ConvertFunc func(fieldValue interface{}, multiplier int)(interface{}, error)

var mapConvertFunc = map[string]ConvertFunc{
	"int": func(fieldValue interface{}, multiplier int)(result interface{}, err error) {
				intValue, err:= strconv.ParseInt(fieldValue.(string), 10, 64)
				if err != nil {
					return nil, err
				}
				if multiplier > 0 {
					intValue *= int64(multiplier)
				}
				return intValue, nil
			},
	"float": func(fieldValue interface{}, multiplier int)(result interface{}, err error) {
				floatValue, err := strconv.ParseFloat(fieldValue.(string), 64)
				if err != nil {
					return nil, err
				}
				if multiplier > 0 {
					floatValue *= float64(multiplier)
				}
				return floatValue, nil
			},
}


func defaultFilterConfig()FilterConfig{
	return FilterConfig{
		CommonConfig: config.CommonConfig{Type:Module},
	}
}


func ConvertFilterHandler(ctx context.Context, raw *config.ConfigRaw)(config.TypeFilterConfig, error){
	conf := defaultFilterConfig()
	if err := config.ReflectConfig(raw, &conf); err != nil {
		return nil, err
	}
	return &conf, nil
}


func (f *FilterConfig)Event(ctx context.Context,event logevent.LogEvent)logevent.LogEvent{
	for _, iRaw := range(f.Fields) {
		field ,err := NewConvertItem(iRaw)
		if err != nil {
			log.Println(err)
			continue
		}

		fieldValue, ok:= event.Event[field.Field]
		if !ok {
			continue
		}
		convertFunc, ok := mapConvertFunc[field.To]
		if !ok {
			log.Println("no convert func type:" + field.To)
			return event
		}
		convertValue, err := convertFunc(fieldValue, field.Multiplier)
		if err != nil {
			log.Println(err)
			if field.RemoveIfFail {
				delete(event.Event, field.Field)
			}
		}
		event.Event[field.Field] = convertValue
	}
	return event
}

