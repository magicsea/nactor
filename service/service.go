package service

import (
	"github.com/magicsea/nactor/actor"
	"github.com/nats-io/nats.go"
)

//服务接口
type IService interface {
	IBaseService
	GetServiceType() string
	actor.ActorProc
}

//开始服务
func StartService(s IService,conn *nats.Conn) error {
	ac := actor.NewActor(s.GetName(),conn,s)
	s.setActor(ac)
	err := ac.Start()
	return err
}
