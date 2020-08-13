package app

import (
	"github.com/spf13/viper"
)

type ServerConf struct {
	App AppConf
	Mq MqConf
	//Log LogConf
	Services []ServiceConf
}
//type LogConf struct {
//	Level string
//	Path string
//	Flag int
//}

type AppConf struct {
	Version string
	Recover bool
}

type MqConf struct {
	Addr string
}

type ServiceConf struct {
	Type   string
	Name string
}

var serverConfig *ServerConf
func LoadServerConfig(fileType string,path string) (*ServerConf,error)  {

	viper.SetConfigType(fileType)
	viper.SetConfigFile(path)
	if err := viper.ReadInConfig(); err != nil {
		return nil,err
	}
	var conf ServerConf
	if err := viper.Unmarshal(&conf);err!=nil {
		return nil,err
	}
	return &conf,nil
}

func GetServerConfig() *ServerConf {
	return serverConfig
}