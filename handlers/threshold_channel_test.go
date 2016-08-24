package handlers

import (
	"testing"
	"bytes"
	"github.com/pbergman/logger"
)

func TestThresholdChannel(t *testing.T) {
	buffer := new(bytes.Buffer)
	record := []logger.Record{
		getRecord("bar", logger.DEBUG, logger.ChannelName("main")),
	}
	handler := NewThresholdChannelHandler(
		NewWriterHandler(buffer, logger.DEBUG),
		map[logger.ChannelName]logger.LogLevel{
			logger.ChannelName("foo"): logger.CRITICAL,
			logger.ChannelName("bar"): logger.ERROR,
		},
		5,
	)
	handler.handler.SetFormatter(&formatter{})

	if true != handler.Support(record[0]) {
		t.Errorf("Expecting to support record %#v", record)
	}
	handler.Handle(&record[0])
	if s := buffer.Len(); s > 0 {
		t.Error("Expecting handler not to be called.")
	}
	for i := 0; i < 10; i++ {
		handler.Handle(&record[0])
	}
	if s := len(handler.buffer); s != 5 {
		t.Errorf("Expecting buffer size not to exceed 5, size: %d, cap %d.", len(handler.buffer), cap(handler.buffer))
	}
	record = append(record, getRecord("foo", logger.ERROR, logger.ChannelName("foo")))
	handler.Handle(&record[1])
	if s := len(handler.buffer); s != 5 {
		t.Errorf("Expecting buffer size not to exceed 5, size: %d, cap %d.", len(handler.buffer), cap(handler.buffer))
	}
	if buffer.Len() > 0 {
		t.Errorf("Logger should not have wirtten got: %s", buffer.String())
	}
	record = append(record, getRecord("foo", logger.ALERT, logger.ChannelName("foo")))
	handler.Handle(&record[2])
	if str := buffer.String(); str != "DEBUG\nDEBUG\nDEBUG\nERROR\nALERT\n" {
		t.Errorf("Expecting: 'DEBUG\nDEBUG\nDEBUG\nERROR\nALERT\n' got: %s size: %d, cap %d", str, len(handler.buffer), cap(handler.buffer))
	}
	buffer.Truncate(0)
	record = append(record, getRecord("foo", logger.DEBUG, logger.ChannelName("main")))
	handler.Handle(&record[3])
	if buffer.Len() != 6 {
		t.Errorf("Expecting: 'DEBUG' got: %s", buffer.String())
	}
	handler.SetStopBuffering(false)
	buffer.Truncate(0)
	record = append(record, getRecord("foo", logger.ERROR, logger.ChannelName("foo")))
	handler.Handle(&record[4])
	record = append(record, getRecord("foo", logger.WARNING, logger.ChannelName("bar")))
	handler.Handle(&record[5])
	if s := len(handler.buffer); s != 2 {
		t.Errorf("Expecting buffer size not to exceed 2, size: %d, cap %d.", len(handler.buffer), cap(handler.buffer))
	}
	record = append(record, getRecord("foo", logger.ERROR, logger.ChannelName("bar")))
	handler.Handle(&record[6])
	if str := buffer.String(); str != "ERROR\nWARNING\nERROR\n" {
		t.Errorf("Expecting: 'ERROR\nWARNING\nERROR\n' got: %s size: %d, cap %d", str, len(handler.buffer), cap(handler.buffer))
	}
}

func TestThresholdChannel_processor(t *testing.T) {
	buffer := new(bytes.Buffer)
	record := getRecord("bar", logger.WARNING, logger.ChannelName("main"))
	handler := NewThresholdChannelHandler(
		NewWriterHandler(buffer, logger.DEBUG),
		map[logger.ChannelName]logger.LogLevel{
			logger.ChannelName("foo"): logger.CRITICAL,
			logger.ChannelName("bar"): logger.ERROR,
		},
		5,
	)
	handler.AddProcessor(func(r *logger.Record) {
		r.Channel = logger.ChannelName("foo")
	})

	if handler.GetProcessors().Len() <= 0 {
		t.Errorf("Expecting to have 1 processor got %d", handler.GetProcessors().Len())
	}

	handler.Handle(&record)

	if record.Channel.GetName() != "foo" {
		t.Errorf("Expecting record to have channel name 'foo' got: %s", record.Channel.GetName())
	}
}

func TestThresholdChannel_channel(t *testing.T) {
	buffer := new(bytes.Buffer)
	record := getRecord("bar", logger.WARNING, logger.ChannelName("main"))
	handler := NewThresholdChannelHandler(
		NewWriterHandler(buffer, logger.DEBUG),
		map[logger.ChannelName]logger.LogLevel{
			logger.ChannelName("foo"): logger.CRITICAL,
			logger.ChannelName("bar"): logger.ERROR,
		},
		5,
		logger.ChannelName("!main"),
	)
	if true == handler.GetChannels().Support(record.Channel) {
		t.Errorf("Handler should not support channel %s (handler: %s)", record.Channel.GetName(), (*handler.channels)[handler.channels.FindChannel("main")])
	}

	if false == handler.GetChannels().Support(logger.ChannelName("test")) {
		t.Errorf("Handler should support channel %s (handler: %s)", record.Channel.GetName(), (*handler.channels)[handler.channels.FindChannel("main")])
	}
}
