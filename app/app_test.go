package app

import (
	"fmt"
	"github.com/magicsea/nactor/actor"
	"github.com/magicsea/nactor/nlog"

	"time"

	"github.com/magicsea/nactor/service"

	"github.com/stretchr/testify/assert"
	"testing"
)

func newExamService() service.IService {
	s := &exampleService{}
	return s
}
type exampleService struct {
	service.BaseService
}
func (s *exampleService) OnInitService() {
	s.RegisterAllRecvMethod(s)
	fmt.Println("OnStart")
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
	nlog.Info("OnRecv_string:",msg)
}

func TestApp(t *testing.T)  {
	conf,err := LoadServerConfig("toml","test.toml")
	assert.NoError(t,err)

	//logger := zaplog.InitLogger()
	//defer logger.Sync()
	//nlog.SetLogger(logger)

	servicename := "exam1"
	RegisterService("exam",newExamService)

	go Run(conf)

	time.Sleep(time.Second)
	for i:=0;i<2;i++  {
		actor.NewProxy(servicename,GetMQConn()).Tell("hello")
	}
	time.Sleep(time.Second)
}
