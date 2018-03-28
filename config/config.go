package config

import (
	"io/ioutil"
	"path/filepath"
	"errors"
	"gopkg.in/yaml.v2"
	"context"
	"golang.org/x/sync/errgroup"
	"github.com/modeyang/LogCruiser/config/logevent"
	"github.com/rcrowley/go-metrics"
	"reflect"
	"syscall"
	"time"
)

// TypeCommonConfig is interface of basic config
type TypeCommonConfig interface {
	GetType() string
}

// CommonConfig is basic config struct
type CommonConfig struct {
	Type string `yaml:"type"`
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
	MetricRaw 	[]ConfigRaw `yaml:"metrics"`
	SinkRaw 	[]ConfigRaw `yaml:"sink"`

	ChanSize 	int64		`yaml:"chsize,omitempty"`
	Interval 	int64 		`yaml:"interval,default 3"`
	SinkTimeRange int64		`yaml:"sinkTimeRange,default 60"`

	chInFilter  MsgChan // channel from input to filter
	chInMetric 	MsgChan // channel from filter to metric
	chInSinker  MetricChan // channel from metric to sink

	registry 		metrics.Registry //metric registry
	selfRegistry 	metrics.Registry // self metric registry

	ctx 			context.Context
	eg 				*errgroup.Group
}

type MsgChan chan logevent.LogEvent
type MetricChan chan MetricResult

var defaultConfig = Config{
	ChanSize: 100,
	Interval: 3,
	SinkTimeRange: 60,
}

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
	if err = yaml.Unmarshal(data, &config); err != nil {
		return config, errors.New("Failed unmarshalling config in YAML format")
	}
	initConfig(&config)
	return
}

func initConfig(config *Config) {
	rv := reflect.ValueOf(&config)
	formatReflect(rv)

	if config.ChanSize < 1 {
		config.ChanSize = defaultConfig.ChanSize
	}

	if config.Interval < 1 {
		config.Interval = defaultConfig.Interval
	}

	if config.SinkTimeRange < 1 {
		config.SinkTimeRange = defaultConfig.SinkTimeRange
	}

	config.chInFilter = make(MsgChan, config.ChanSize)
	config.chInMetric = make(MsgChan, config.ChanSize)
	config.chInSinker = make(MetricChan, defaultConfig.ChanSize)

	config.registry = metrics.NewRegistry()
	config.selfRegistry = metrics.NewRegistry()

}

func (c *Config) handleMetrics()error{
	rawMetrics := map[string]int64{}
	allMetrics := c.registry.GetAll()
	c.registry.UnregisterAll()

	for k, v := range(allMetrics) {
		if c, ok := v["count"]; ok {
			rawMetrics[k] = c.(int64)
		}
		if c, ok := v["value"]; ok {
			rawMetrics[k] = c.(int64)
		}
	}
	now := time.Now().Second()
	timestamp := now - now % int(c.SinkTimeRange)
	c.chInSinker <- MetricResult{Timestamp: timestamp, Data: rawMetrics}
	return nil
}

func (c *Config)Start(ctx context.Context) (err error){
	ctx = contextWithOSSignal(ctx, syscall.SIGINT, syscall.SIGTERM)
	c.eg, c.ctx = errgroup.WithContext(ctx)

	if err = c.startInput(); err != nil {
		return
	}
	if err = c.startFilters(); err != nil {
		return
	}
	if err = c.startMetrics(); err != nil {
		return
	}
	if err = c.startSinkers(); err != nil {
		return
	}

	// new ticker with interval for sink metrics
	ticker := time.NewTicker(time.Second * time.Duration(c.Interval))
	go func() {
		for _ = range(ticker.C) {
			err = c.handleMetrics()
		}
	}()
	return
}

func (c *Config)Wait()error{
	return c.eg.Wait()
}



