package service

import (
	"fmt"
	"github.com/magicsea/nactor/actor"
	"github.com/magicsea/nactor/encode"
	"github.com/nats-io/nats.go"

	"testing"
	"time"
)

type exampleService struct {
	BaseService
}
func (s *exampleService) OnInitService() {
	fmt.Println("OnInit")
	s.RegisterAllRecvMethod(s)
}
func (s *exampleService) OnStartService() {
	fmt.Println("OnStart")
}
func (s *exampleService) OnDestroy() {
	fmt.Println("OnDestroy")
}

func (s *exampleService) GetServiceType() string {
	return "exm"
}
func (s *exampleService) OnRecv_string(ctx actor.Context,msg string)  {
	fmt.Println("OnRecv_string:",msg)
}
func (s *exampleService) OnRecv_St(ctx actor.Context,msg *St)  {
	fmt.Println("OnRecv_St:",msg)
}

type St struct {
	A int
	S string
}
func TestRegFunc(t *testing.T)  {
	encode.RegisterName((*St)(nil))

	nc, _ := nats.Connect(nats.DefaultURL)
	sv := &exampleService{}
	sv.Init("exm1")

	StartService(sv,nc)
	defer sv.Stop()

	go sv.actor.Run()

	//sv.actor.Read()//actor.Started

	time.Sleep(time.Second)
	s := "hello"
	actor.NewProxy(sv.GetName(),nc).Tell(s)
	//assert.EqualValues(t,sv.actor.Read().Message(),s)

	st := &St{A:1,S:"ffff"}
	actor.NewProxy(sv.GetName(),nc).Tell(st)
	//assert.EqualValues(t,sv.actor.Read().Message(),st)

	time.Sleep(time.Second)
}