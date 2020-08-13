package app

import (
	"fmt"
	"github.com/magicsea/nactor/module"
	"github.com/magicsea/nactor/nlog"
	"github.com/magicsea/nactor/service"
	"github.com/magicsea/nactor/util"
	"github.com/nats-io/nats.go"

	"os"
	"os/signal"
	"reflect"
)

type MakeServiceFunc func() service.IService

var (
	serviceTypeMap map[string]MakeServiceFunc
	services       []service.IService
	modules        []module.IModule
	mqConn 			*nats.Conn
)

func init() {
	serviceTypeMap = make(map[string]MakeServiceFunc)
}

//mq连接
func GetMQConn() *nats.Conn {
	return mqConn
}

func RegisterService(serviceType string, f MakeServiceFunc) {
	serviceTypeMap[serviceType] = f
}

/**
 *启动app
 *会阻塞在此状态，直到手动关闭
 *
**/
func Run(conf *ServerConf,ms ...module.IModule) {
	if conf.App.Recover {
		defer util.PrintPanicStack()
	}
	serverConfig = conf

	nlog.Info("log started.")
	//modules init
	modules = ms
	for _, m := range modules {
		if !m.OnInit() {
			panic(fmt.Sprintf("%v module.OnInit fail", m))
		}
	}
	for _, m := range modules {
		m.Run()
	}

	//start mq
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		panic(err)
	}
	mqConn = nc

	//service init
	for _, sc := range conf.Services {
		makefunc := serviceTypeMap[sc.Type]
		if makefunc != nil {
			ser := makefunc()
			nlog.Infof("create service:%s", sc.Name)
			ser.Init(sc.Name)
			services = append(services, ser)
		} else {
			panic(fmt.Sprintf("not found service type:%v", sc))
		}
	}

	//start
	for _, ser := range services {
		nlog.Infof("start service:%s", ser.GetName())
		if err:=service.StartService(ser,mqConn);err!=nil {
			panic(err)
		}
	}
	nlog.Info("all service started!")

	// close
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	sig := <-c
	nlog.Infof("closing down (signal: %v)", sig)
	OnDestroy()
}

func OnDestroy() {
	for _, ser := range services {
		nlog.Infof("destroy service:%s", ser.GetName())
		ser.OnDestroy()
	}
	for _, m := range modules {
		nlog.Infof("destroy module:%v", reflect.TypeOf(m))
		m.OnDestroy()
	}
}
