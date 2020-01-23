package async

import (
	"time"
)

type HandlerFunc func() (interface{}, error)

type Async struct {
	actions []Action
}

func (a *Async) Add(name string, handler HandlerFunc) {
	a.actions = append(a.actions, Action{
		name:    name,
		handler: handler,
	})
}

func (a *Async) Process() []Response {
	if len(a.actions) == 0 {
		return nil
	}

	r := []Response{}
	ch := make(chan Response)

	for _, a := range a.actions {
		go func(a Action) {
			v, err := a.handler()
			ch <- Response{
				Name:  a.name,
				Value: v,
				Error: err,
			}
		}(a)
	}

	for {
		select {
		case resp := <-ch:
			r = append(r, resp)

			if len(r) == len(a.actions) {
				return r
			}
		default:
			time.Sleep(time.Millisecond * 50)
		}
	}

	return nil
}

type Action struct {
	name    string
	handler HandlerFunc
}

type Response struct {
	Name  string
	Value interface{}
	Error error
}
