package logger

import (
	"bytes"
	"fmt"
	"log"

	"github.com/sirupsen/logrus"
)

type CustomFormatter struct {
}

func (f CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	keys := make([]string, 0, len(entry.Data))
	for k := range entry.Data {
		keys = append(keys, k)
	}

	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	prefixFieldClashes(entry.Data)

	f.appendKeyValue(b, "level", entry.Level.String())
	if entry.Message != "" {
		f.appendKeyValue(b, "msg", entry.Message)
	}
	for _, key := range keys {
		f.appendKeyValue(b, key, entry.Data[key])
	}

	b.WriteByte('\n') //nolint: gosec
	return b.Bytes(), nil
}

//nolint: gosec
func (f *CustomFormatter) appendKeyValue(b *bytes.Buffer, key string, value interface{}) {

	b.WriteByte('[')
	b.WriteString(key)
	b.WriteByte('=')
	f.appendValue(b, value)
	b.WriteByte(']')
	b.WriteByte(' ')
}

//nolint: gosec
func (f *CustomFormatter) appendValue(b *bytes.Buffer, value interface{}) {
	switch value := value.(type) {
	case string:
		b.WriteString(value)
	case error:
		errmsg := value.Error()
		b.WriteString(errmsg)
	default:
		_, err := fmt.Fprint(b, value)
		if err != nil {
			log.Println("Unable to log: ", value)
		}
	}
}

func prefixFieldClashes(data logrus.Fields) {
	if t, ok := data["time"]; ok {
		data["fields.time"] = t
	}

	if m, ok := data["msg"]; ok {
		data["fields.msg"] = m
	}

	if l, ok := data["level"]; ok {
		data["fields.level"] = l
	}
}
