package logger

import (
	"io"
)

// NewThresholdHandler will return an handler that buffers X amount of (size) records
// till the level is reached on which it will send all records the wrapped handler
func NewThresholdHandler(handler HandlerInterface, size int, level LogLevel, bubble bool) HandlerInterface {
	return &thresholdHandler{buffer: newRecordBuffer(size), handler: handler, baseHandler: baseHandler{level: level, bubble: bubble}}
}

type thresholdHandler struct {
	buffer  *recordBuffer
	handler HandlerInterface
	baseHandler
	processor
}

func (t thresholdHandler) IsHandling(r *Record) bool {
	return true
}

func (t *thresholdHandler) Handle(r *Record) bool {
	t.process(r)
	t.buffer.push(r)
	if r.Level <= t.level {
		for t.buffer.valid() {
			t.handler.Handle(t.buffer.shift())
		}
	}
	return t.bubble
}

func (f *thresholdHandler) Close() error {
	var err error
	if closer, ok := f.handler.(io.Closer); ok {
		err = closer.Close()
	}
	return err
}
