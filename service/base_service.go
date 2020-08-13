package service

import (
	"errors"
	"github.com/magicsea/nactor/actor"
	"github.com/magicsea/nactor/nlog"
	"reflect"
	"regexp"
)
//基础服务接口
type IBaseService interface {
	Init(name string)
	GetName() string

	setActor(actor actor.IActor)
}

//基础服务
type BaseService struct {
	name     string
	actor actor.IActor
	rounter    map[reflect.Type]reflect.Value
}

func (s *BaseService) Init(name string) {
	s.rounter = make(map[reflect.Type]reflect.Value)
	s.name = name
}

func (s *BaseService) GetName() string {
	return s.name
}
func  (s *BaseService) setActor(actor actor.IActor) {
	s.actor = actor
}
func  (s *BaseService) Stop()  {
	s.actor.Close()
}

func (s *BaseService) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *actor.Started:
		nlog.Info("Started, initialize actor here:",s.name)

	default:
		nlog.Debug("service recv defalult:",s.name,"=>", msg)
		err := s.CallMethod(context)
		if err != nil {
			nlog.Error(err)
		    return
		}
	}
}

/**
 *注册一个接受消息函数
 *函数格式:
 *func(arg pb.Hello,ctx actor.Context)
 *arg参数只能有一个
**/
func (s *BaseService) RegisterMsg(f interface{}) {
	s.rounter[reflect.TypeOf(f)] = reflect.ValueOf(f)
}

/**
 *注册对象的所有接受消息函数
 *函数名无要求，参数要求(actor.Context,消息类型)
 *e.g. OnRecv_string(ctx actor.Context,msg string)
 *
**/
func (s *BaseService) RegisterAllRecvMethod(rawPtr interface{}) error {
	ptrValue := reflect.ValueOf(rawPtr)

	for i :=1;i<ptrValue.NumMethod() ;i++  {
		m := ptrValue.Method(i)
		methodName := m.String()
		if match,_ := regexp.MatchString("<func\\(actor.Context,(.*)\\) Value>",methodName);match {
			tp := reflect.TypeOf(rawPtr).Method(i)
			msgType :=  tp.Type.In(2)//第一个参数是消息体
			if !m.IsValid() {
				return errors.New("RegisterAllRecvMethod method error:"+methodName)
			}
			s.rounter[msgType] = m
		}
	}
	return nil
}

//CallMethod
func (s *BaseService) CallMethod(ctx actor.Context) error {
	msg := ctx.Message()
	m,ok := s.rounter[reflect.TypeOf(msg)]
	if !ok {
		return errors.New("no method found:"+reflect.TypeOf(msg).String())
	}

	argValues := make([]reflect.Value, 0, 2)
	argValues = append(argValues, reflect.ValueOf(ctx))
	argValues = append(argValues, reflect.ValueOf(msg))
	m.Call(argValues)
	return nil
}

