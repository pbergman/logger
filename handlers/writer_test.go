package handlers

import (
	"time"
	"testing"
	"bytes"
	"reflect"
	"github.com/pbergman/logger"
	"github.com/pbergman/logger/formatters"
)

var test_time time.Time = time.Date(2016, 1, 2, 10, 20, 30, 0, time.Local)

func getRecord(m string, l logger.LogLevel, n logger.ChannelName) logger.Record {
	return logger.Record{
		Time:    test_time,
		Channel: n,
		Context: make(map[string]interface{}),
		Message: m,
		Level:   l,
	}
}

func TestWriter(t *testing.T) {
	record := getRecord("bar", logger.WARNING, logger.ChannelName("main"))
	buff := new(bytes.Buffer)
	handler := NewWriterHandler("foo", buff, logger.INFO)
	if false == handler.Support(record) {
		t.Errorf("Handler should support level %s (%s > %s)", record.Level, record.Level, handler.level)
	}
	if true == handler.Support(getRecord("bar", logger.DEBUG, logger.ChannelName("main"))) {
		t.Errorf("Handler should not support level %s (%s < %s)", logger.DEBUG, logger.DEBUG, handler.level)
	}
	if err := handler.GetFormatter().Format(record, handler); err != nil {
		t.Error(err.Error())
	}
	if buff.String() != "[2016-01-02 10:20:30.000000] main.WARNING: bar \n" {
		t.Errorf("Expecting '[2016-01-02 10:20:30.000000] main.WARNING: bar' got: %s", buff.String())
	}
}

func TestWriter_channel(t *testing.T) {
	buff := new(bytes.Buffer)
	record := getRecord("bar", logger.WARNING, logger.ChannelName("main"))
	handler := NewWriterHandler("foo", buff, logger.INFO)
	defer handler.Close()
	handler.GetChannels().AddChannel(logger.ChannelName("!main"))
	if true == handler.Support(record) {
		t.Errorf("Handler should not support channel %s (handler: %s)", record.Channel.GetName(), (*handler.channels)[handler.channels.FindChannel("main")])
	}
	record.Channel = logger.ChannelName("not_exist")
	if true == handler.Support(record) {
		t.Errorf("Handler should not support channel %s", record.Channel.GetName())
	}
}

type testWriter struct {}
func (w testWriter) Write(p []byte) (n int, err error) {return 0, nil }
func (w testWriter) Close() error { return nil }

func TestWriter_channels_not_nil(t *testing.T) {
	handler := &WriterHandler{}
	if handler.GetChannels() == nil {
		t.Errorf("Expecting channels not to be nil")
	}
}

func TestWriter_channe(t *testing.T) {
	handler := NewWriterHandler(
		"foo",
		&WriterHandler{},
		logger.INFO,
		logger.ChannelName("foo"),
		logger.ChannelName("bar"),
		logger.ChannelName("example"),
	)

	if handler.GetChannels().Len() != 3 {
		t.Errorf("Expecting 3 channels got %d", handler.GetChannels().Len() )
	}
}

func TestWriter_getters(t *testing.T) {
	writer := testWriter{}
	handler := NewWriterHandler("foo", writer, logger.INFO)
	defer handler.Close()

	if handler.GetName() != "foo" {
		t.Errorf("Expecting name 'foo' got:  %s", handler.GetName())
	}

	if handler.GetLevel() != logger.INFO {
		t.Errorf("Expecting level 'INFO' got:  %s", logger.INFO)
	}

	formatter := formatters.NewCustomLineFormatter("{{ .Message }}")
	handler.SetFormatter(formatter)

	if !reflect.DeepEqual(formatter, handler.GetFormatter()) {
		t.Errorf("Expecting to get same formatter as was set (%p != %p)", formatter, handler.GetFormatter())
	}

	channels := new(logger.ChannelNames)
	channels.AddChannel(logger.ChannelName("foo"))
	channels.AddChannel(logger.ChannelName("bar"))
	handler.SetChannels(channels)

	if !reflect.DeepEqual(channels, handler.GetChannels()) {
		t.Errorf("Expecting to get same channel colection as was set (%p != %p)", formatter, handler.GetFormatter())
	}
}