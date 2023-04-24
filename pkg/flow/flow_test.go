package flow

import (
	"dh/pkg/config"
	"errors"
	"testing"
)

func TestFlow_ExecuteFlow(t *testing.T) {
	t.Run("Flow with two actions", func(t *testing.T) {
		var flowHandlers = []*EffectHandler{
			{
				Action: func() error {
					config.InfoLog.Print("action 1")
					return nil
					//return errors.New("1 err")
				},
				OnError: func(err error) {
					config.WarnLog.Print("action 1 onError")
				},
			},
			{
				Action: func() error {
					config.InfoLog.Print("action 2")
					//return nil
					return errors.New("2 err")
				},
				OnError: func(err error) {
					config.WarnLog.Print("action 2 onError")
				},
			},
		}
		f := &Flow{}
		_ = f.CreateFlow(flowHandlers...)
		err := f.ExecuteFlow()
		if err == nil {
			t.Fail()
		}
	})
}
