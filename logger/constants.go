package logger

type LogField string

const (
	Application   LogField = "application"
	Platform      LogField = "platform"
	Environment   LogField = "environment"
	AmazonTraceID LogField = "x-amzn-trace-id"
	Error         LogField = "error"
	Header        LogField = "header"
	Url           LogField = "url"
	Uri           LogField = "uri"
	Hostname      LogField = "hostname"
	Body          LogField = "body"
	Status        LogField = "status"
	Method        LogField = "cs-method"
	Path          LogField = "path"
	Pattern       LogField = "cs-url"
	StackTrace    LogField = "stacktrace"
	StatusCode    LogField = "cs-status"
	TimeTaken     LogField = "time-taken-ms"
	ResponseBody  LogField = "response_body"
)
