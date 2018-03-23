package config

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
}

//func LoadFromFile(path string)(config Config, err error) {
//
//}
//
//func LoadFromYaml(data []byte)(config Config, err error){
//
//}