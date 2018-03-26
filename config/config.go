package config

import (
	"io/ioutil"
	"path/filepath"
	"errors"
	"gopkg.in/yaml.v2"
	"log"
	"context"
	"golang.org/x/sync/errgroup"
	"github.com/modeyang/LogCruiser/config/logevent"
	"github.com/rcrowley/go-metrics"
)

// TypeCommonConfig is interface of basic config
type TypeCommonConfig interface {
	GetType() string
}

// CommonConfig is basic config struct
type CommonConfig struct {
	Type string `json:"type"`
}

// GetType return module type of config
func (t CommonConfig) GetType() string {
	return t.Type
}

// ConfigRaw is general config struct
type ConfigRaw map[string]interface{}


// yaml config for all
type Config struct {
	InputRaw 	[]ConfigRaw `yaml:"input"`
	FilterRaw 	[]ConfigRaw `yaml:"filter"`
	MetricRaw 	[]ConfigRaw `yaml:"metric"`
	SinkRaw 	[]ConfigRaw `yaml:"sink"`

	ChanSize 	int			`yaml:"chsize,omitempty"`
	Interval 	int 		`yaml:"interval,default 3"`

	chInFilter  MsgChan // channel from input to filter
	chInMetric 	MsgChan // channel from filter to metric

	registry 		*metrics.StandardRegistry  	//metric registry
	selfRegistry 	*metrics.StandardRegistry	// self metric registry

	ctx 			context.Context
	eg 				*errgroup.Group
}

type MsgChan chan logevent.LogEvent


func LoadFromFile(path string)(config Config, err error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	ext := filepath.Ext(path)
	if ext != ".yml" && ext != ".yaml" {
		return Config{}, errors.New("config file need yaml config")
	}
	return LoadFromYaml(data)
}

func LoadFromYaml(data []byte)(config Config, err error){
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		panic(err)
	}
	log.Println(config.InputRaw)
	return config, nil
}

func (c *Config)start(ctx context.Context) {

}



