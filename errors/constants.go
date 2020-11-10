package errors

const (
	ERROR_READING_REQUEST_BODY     string = "Error reading request body"
	ERROR_MARSHALLING_REQUEST_BODY string = "Error decoding request body: %s"
	REQUEST_BODY_MISSING           string = "Please send a request body"
	PROCESSING_ERROR               string = "PROCESSING_ERROR"
	PROCESSING_SUCCESS             string = "PROCESSING_SUCCESS"
)
