package zaplog

import (
	"log"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"github.com/natefinch/lumberjack"
)

var sugarLogger *zap.SugaredLogger

//InitLogger
func InitLogger() *zap.SugaredLogger{
	path := viper.GetString("log.path")
	slevel := viper.GetString("log.level")
	var level zapcore.Level
	if err:=level.UnmarshalText([]byte(slevel));err!=nil {
		log.Println("log.level is invalid:",slevel)
		panic(err)
	}

	writeSyncer := getLogWriter(path)
	encoder := getEncoder(true)

	core := zapcore.NewCore(encoder, writeSyncer, zapcore.Level(level))

	logger := zap.New(core, zap.AddCaller())
	sugarLogger = logger.Sugar()
	return sugarLogger
}

func getEncoder(isJson bool) zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig() //发布配置
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder //标准时间
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	if isJson {
		zapcore.NewJSONEncoder(encoderConfig)
	}
	return zapcore.NewJSONEncoder(encoderConfig)
}

func getLogWriter(path string) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   path,
		MaxSize:    1,//日志分割大小MB
		MaxBackups: 5,//备份数量
		MaxAge:     30,//保留日期
		Compress:   false,
	}
	return zapcore.AddSync(lumberJackLogger)
}
