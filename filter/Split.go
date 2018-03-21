package filter

import (
	"strings"
)

type SplitFilter struct {
	Source	string  `yaml:"source:default message"`
	Fields []string `yaml:"fields"`
	Split 	string 	`yaml:"split"`
	RemoveFields []string `yaml:"remove_fields:omitempty"`
	ErrorTag string
}

func NewSplitFilter(source string, fields []string, split string, remove_fields []string)*SplitFilter {
	if source == "" {
		source = "message"
	}
	return &SplitFilter{source, fields, split, remove_fields, "splitfailed"}
}

func (f *SplitFilter) Filter(event map[string]interface{})(map[string]interface{}, error) {
	if value, ok := event[f.Source]; ok {
		msg_list := strings.Split(value.(string), f.Split)
		for i, f := range(f.Fields) {
			if i <= len(msg_list) - 1 {
				event[f] = string(msg_list[i])
			} else {
				event[f] = ""
			}
		}
		if len(f.RemoveFields) > 0 {
			for _, k := range(f.RemoveFields) {
				delete(event, k)
			}
		}
	}
	return event, nil
}
