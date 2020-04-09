package logger

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
)

var log = logrus.New()
var CurDir string = ""

func init() {

	rolllogger := &lumberjack.Logger{
		// 日志输出文件路径
		Filename:   "logs/liuyi-test.log",
		// 日志文件最大 size, 单位是 MB
		MaxSize:    20, // megabytes
		// 最大过期日志保留的个数
		MaxBackups: 1000,
		// 保留过期文件的最大时间间隔,单位是天
		MaxAge:     30,   //days
		// 是否需要压缩滚动日志, 使用的 gzip 压缩
		Compress:   true, // disabled by default
	}


	log.SetOutput(io.MultiWriter(rolllogger,os.Stdout))

	//设置输出样式，自带的只有两种样式logrus.JSONFormatter{}和logrus.TextFormatter{}
	log.SetFormatter(&logrus.JSONFormatter{})
	//设置output,默认为stderr,可以为任何io.Writer，比如文件*os.File
	//设置最低loglevel
	log.SetLevel(logrus.InfoLevel)
}

func NewLogger(traceId string,spanId string) *logrus.Entry{
	entry := log.WithFields(logrus.Fields{"trace": traceId, "span": spanId})
	return entry

}
