package file

import (
	"github.com/modeyang/LogCruiser/config"
	"context"
	"github.com/modeyang/LogCruiser/config/logevent"
	"github.com/hpcloud/tail"
	"strings"
	"errors"
)

type InputConfig struct {
	config.InputConfig

	Path 		string  `yaml:"path"`
	// file position, 0 -> beginning  1 -> cur  2 -> end
	Position 	int 	`yaml:"position"`
}

const Module = "file"

func defaultFileInputConfig()(InputConfig){
	return InputConfig{
		InputConfig: config.InputConfig{
			CommonConfig: config.CommonConfig{
				Type:Module,
			},
		},
		Position: 2,
	}
}

func FileInputHandler(ctx context.Context, raw *config.ConfigRaw)(config.TypeInputConfig, error) {
	conf := defaultFileInputConfig()
	err := config.ReflectConfig(raw, &conf)
	if err != nil {
		return nil, err
	}
	if conf.Position > 2 || conf.Position < 0 {
		panic(errors.New("position must in [0, 1, 2]"))
	}
	return &conf, nil
}

func (input *InputConfig)Source(ctx context.Context, msgChan chan <- logevent.LogEvent)(err error){
	location := &tail.SeekInfo{Offset:0, Whence:input.Position}
	t, err := tail.TailFile(input.Path, tail.Config{Follow:true, Location:location})
	if err != nil {
		t.Stop()
		return err
	}
	defer t.Cleanup()
	for {
		select {
		case <- ctx.Done():
			return nil
		case line := <- t.Lines:
			if line != nil && len(line.Text) > 0 {
				// for windows line "\r\n"
				msgChan <- *logevent.NewLogEvent(strings.TrimRight(line.Text, "\r"))
			}
		}
	}
	t.Wait()
	return nil
}

