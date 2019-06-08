package simfaas

import (
	"errors"
)

const (
	CUSTOM_FN_HEADER = "X-CustomFn"
)

type CustomFn = func(b []byte) ([]byte, error)

type CustomHandler struct {
	// a parser function to extract a key used look up a custom handler in
	// the map of available handlers
	GetFn    func(name string) (string, error)
	Handlers map[string]CustomFn
}

func (c CustomHandler) ExecFn(fnName string, b []byte) ([]byte, error) {
	customHandler, err := c.GetFn(fnName)
	if err != nil {
		return nil, err
	}

	fn, ok := c.Handlers[customHandler]
	if !ok {
		return nil, errors.New("Parsed handler does not exist")
	}

	// execute and return report
	return fn(b)
}
