package service

import (
	"github.com/magicsea/nactor/actor"
	"github.com/nats-io/nats.go"
)

//服务接口
type IService interface {
	IBaseService
	//运行在main线程
	OnInitService()
	OnStartService()
	GetServiceType() string
	actor.ActorProc
}

//开始服务
func StartService(s IService,conn *nats.Conn) error {
	ac := actor.NewActor(s.GetName(),conn,s)
	err := ac.Start()
	if err != nil {
	    return err
	}
	s.setActor(ac)
	s.OnInitService()
	go RunService(s,ac)
	return nil
}

//启动服务actor线程
func RunService(s IService,iActor actor.IActor)  {
	s.OnStartService()
	iActor.Run()
}
