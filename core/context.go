package core

import "github.com/sirupsen/logrus"

type TraceContext struct {
	traceConfig TraceConfig
	data map[string] interface{}
	http HttpConfig
	socket SocketConfig
}

type SpanContext struct {
	spanConfig SpanConfig
	params map[string]interface{}
	headers map[string]string
	url string
	result string
	status int
	log *logrus.Entry


}




