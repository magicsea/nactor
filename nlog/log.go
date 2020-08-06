package nlog

import "log"

type Logger interface {
	Debug(v ...interface{})
	Info(v ...interface{})
	Error(v ...interface{})
}

var logger Logger = &defaultLogger{}

//设置日志管理器
func SetLogger(l Logger)  {
	logger = l
}

func Debug(v ...interface{}) {
	logger.Debug(v...)
}
func Info(v ...interface{}) {
	logger.Info(v...)
}
func Error(v ...interface{}) {
	logger.Error(v...)
}


type defaultLogger struct {
}

func (l *defaultLogger) Debug(v ...interface{}) {
	log.Println(v...)
}
func (l *defaultLogger) Info(v ...interface{}) {
	log.Println(v...)
}
func (l *defaultLogger) Error(v ...interface{}) {
	log.Println(v...)
	log.Fatal(v...)
}
