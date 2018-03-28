package logevent

import (
	"time"
	"bytes"
	"log"
	"text/template"
)

type LogEvent struct {
	Timestamp 	time.Time
	Message  	string
	Tags 		[]string
	Event 		map[string]interface{}
}


const timeFormat string = `2006-01-02T15:04:05.999999999Z`

func NewLogEvent(msg string)*LogEvent {
	return &LogEvent{
		Timestamp: time.Now(),
		Message: msg,
		Tags: []string{},
		Event:make(map[string]interface{}),
	}
}

func (t *LogEvent)AddTag(tags ...string) {
	for _, tag := range tags {
		t.Tags = append(t.Tags, tag)
	}
}

