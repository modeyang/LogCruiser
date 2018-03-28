package config

import (
	"reflect"
	"context"
	"os"
	"os/signal"
	"log"
	"bytes"
	"text/template"
	"encoding/json"
	"github.com/icza/dyno"
)

func contextWithOSSignal(parent context.Context, sig ...os.Signal) context.Context {
	osSignalChan := make(chan os.Signal, 1)
	signal.Notify(osSignalChan, sig...)

	ctx, cancel := context.WithCancel(parent)

	go func(cancel context.CancelFunc) {
		select {
		case sig := <-osSignalChan:
			log.Println(sig)
			cancel()
		}
	}(cancel)

	return ctx
}

// ReflectConfig set conf from confraw
func ReflectConfig(confraw *ConfigRaw, conf interface{}) (err error) {
	data, err := json.Marshal(dyno.ConvertMapI2MapS(map[string]interface{}(*confraw)))
	if err != nil {
		return
	}

	if err = json.Unmarshal(data, conf); err != nil {
		return
	}

	rv := reflect.ValueOf(conf).Elem()
	formatReflect(rv)

	return
}

func formatReflect(rv reflect.Value) {
	if !rv.IsValid() {
		return
	}

	switch rv.Kind() {
	case reflect.Ptr:
		if !rv.IsNil() {
			formatReflect(rv.Elem())
		}
	case reflect.Struct:
		for i := 0; i < rv.NumField(); i++ {
			field := rv.Field(i)
			formatReflect(field)
		}
	case reflect.String:
		if !rv.CanSet() {
			return
		}
		value := rv.Interface().(string)
		rv.SetString(value)
	}
}

func RenderTemplate(tmpl *template.Template, event map[string]interface{}) (string, error) {
	var tpl bytes.Buffer
	err := tmpl.Execute(&tpl, event)
	if err != nil {
		log.Fatalln(err)
		return "", err
	}
	return tpl.String(), nil
}

func Map2String(m map[string]interface{}) string{
	s, _:= json.Marshal(m)
	return string(s)
}
