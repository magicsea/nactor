package nlog

import (
	"fmt"
	"log"
)


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
func (l *defaultLogger) Debugf(info string,v ...interface{}) {
	log.Printf(info,v...)
}
func (l *defaultLogger) Infof(info string,v ...interface{}) {
	log.Printf(info,v...)
}
func (l *defaultLogger) Errorf(info string,v ...interface{}) {
	log.Printf(info,v...)
	log.Fatal(fmt.Sprintf(info,v...))
}

