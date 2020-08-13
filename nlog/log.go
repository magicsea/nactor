package nlog


type Logger interface {
	Debug(v ...interface{})
	Info(v ...interface{})
	Error(v ...interface{})
	Debugf(info string,v ...interface{})
	Infof(info string,v ...interface{})
	Errorf(info string,v ...interface{})
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
func Debugf(info string,v ...interface{}) {
	logger.Debugf(info,v...)
}
func Infof(info string,v ...interface{}) {
	logger.Infof(info,v...)
}
func Errorf(info string,v ...interface{}) {
	logger.Errorf(info,v...)
}
