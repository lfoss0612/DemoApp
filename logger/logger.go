package logger

import (
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	demoerrors "github.com/lfoss0612/DemoApp/errors"
)

type Logger struct {
	*logrus.Entry
}

var logger = logrus.New()

//InitLogger initialize looger
func Init(logLevel string, out io.Writer, formatter logrus.Formatter) {

	logger.Formatter = formatter

	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		log.Fatal(err.Error())
	}

	logger.SetLevel(level)

	logger.Out = out

	log.SetOutput(logger.Writer())
}

func NewLogger() *Logger {
	return &Logger{
		Entry: logrus.NewEntry(logger),
	}
}

func NewStdLogger() logrus.FieldLogger {
	return logrus.NewEntry(logger)
}

func GetLogger() *logrus.Logger {
	return logger
}

func GetCause(err error) string {
	if err != nil {
		cause := errors.Cause(err)
		if cause == nil {
			cause = err
		}
		return cause.Error()
	}
	return ""
}

func GetStack(err error) []string {
	stack := make([]string, 0)
	if err, ok := err.(demoerrors.StackTracer); ok {
		for _, f := range err.StackTrace() {
			frame := ParseStackTrace(fmt.Sprintf("%+s:%d", f, f))
			stack = append(stack, frame...)
		}
	}
	return stack
}

func ParseStackTrace(stack string) []string {
	stackArr := strings.Split(stack, "\n")

	newStackArr := make([]string, 0)
	for _, v := range stackArr {
		newString := strings.Replace(v, "\t", "     ", -1)
		if newString != "" {
			newStackArr = append(newStackArr, newString)
		}
	}
	return newStackArr
}

func (logger *Logger) IsDebugEnabled() bool {
	return logger.Logger.Level == logrus.DebugLevel
}

func (logger *Logger) AddConstantField(key LogField, value interface{}) {
	logger.Entry = logger.Entry.WithField(string(key), value)
}

func (logger *Logger) WithLogField(key LogField, value interface{}) *Logger {
	logger.Entry = logger.Entry.WithField(string(key), value)
	return logger
}

func (logger *Logger) WithLogFields(fieldMap map[LogField]interface{}) *Logger {
	fields := make(logrus.Fields)
	for k, v := range fieldMap {
		fields[string(k)] = v
	}
	logger.Entry = logger.Entry.WithFields(fields)
	return logger
}

func (logger *Logger) Error(args ...interface{}) {
	if len(args) == 1 {
		if err, ok := args[0].(error); ok {
			logger.WithLogField(StackTrace, GetStack(err)).Error(err.Error())
		}
	} else {
		logger.Entry.Error(args...)
	}
}
