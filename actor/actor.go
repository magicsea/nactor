package actor

import (
	"fmt"
	"github.com/magicsea/nactor/nlog"
	"github.com/magicsea/nactor/encode"
	"github.com/nats-io/nats.go"
	"golang.org/x/net/context"
	"errors"
	"sync"
)
//size of recv chan
const MsgRecvSize = 256

//Actor
type Actor struct {
	name string
	conn *nats.Conn
	ch chan Context
	ctxCancel context.Context
	cancel context.CancelFunc
	subjects sync.Map
	proc ActorProc
}

//NewActor
func NewActor(name string,conn *nats.Conn,proc ActorProc) *Actor {
	ac := Actor{name:name,conn:conn,proc:proc}
	return &ac
}

//start actor
func (ac *Actor) Start() error {
	subject := ac.subjectName()
	ch := make(chan Context,MsgRecvSize)
	ac.ch = ch
	err := ac.Subscribe(subject)
	ac.ctxCancel,ac.cancel = context.WithCancel(context.Background())
	if ac.proc!=nil {
		ac.proc.OnStart()
	}
	return err
}

//发布信息到此actor主题,goroutine safe
func (ac *Actor) Tell(message Message) error {
	data,err:= encode.Encode(message)
	if err != nil {
	    return err
	}
	return ac.conn.Publish(ac.subjectName(),data)
}

/**
 *订阅一个主题,goroutine safe
 *一个actor默认会订阅一个自己名字的主题。
 *上层逻辑可以自己订阅关心的额外主题。
**/
func (ac *Actor) Subscribe(subject string) error {
	nc := ac.conn
	s,errsub := nc.Subscribe(subject, func(m *nats.Msg) {
		nlog.Debug(fmt.Sprintf("[%s] Received a message:%s=> %s",ac.name,subject, string(m.Data)))
		//反序列化
		msg,err := encode.Decode(m.Data)
		if err!=nil {
			nlog.Error(err)
		}
		ctx := newContext(m,msg)
		ac.ch<-ctx
	})

	if errsub != nil {
		nlog.Error("sub error:",errsub)
		return errsub
	}
	ac.subjects.Store(subject,s)
	return nil
}

/**
 *释放一个主题,goroutine safe
 *额外订阅的主题需要上层主动释放，否则会等actor销毁一起释放
 *
**/
func (ac *Actor) Unsubscribe(subject string) error {
	v,ok := ac.subjects.Load(subject)
	if !ok {
		return errors.New("not found subject")
	}
	v.(*nats.Subscription).Unsubscribe()
	ac.subjects.Delete(subject)

	return nil
}

/**
 *主生命期
 *
 *主要职责是读消息，控制生命期，释放资源
**/
func (ac *Actor) Run() error {
	//ec, err := nats.NewEncodedConn(ac.conn, nats.JSON_ENCODER)
	//if err != nil {
	//    return err
	//}
	//defer ec.Close()
	ctx := ac.ctxCancel
	for {
		select {
		case <-ctx.Done():
			goto BREAK
		case c := <-ac.ch:
			nlog.Debug(fmt.Sprintf("[%s] do a message: %v",ac.name, c.Message()))
			if ac.proc!=nil {
				ac.proc.Receive(c)
			}
		}
	}
BREAK:
	ac.onDestroy()
	return nil
}

/**
 *阻塞读一条消息
 *不使用Run()接管消息，自己管理消息接收时候使用
 *
**/
func  (ac *Actor) Read() Context {
	c := <-ac.ch
	return c
}


func (ac *Actor) subjectName() string {
	//aname := fmt.Sprintf("__actor#%s",ac.name)
	return ac.name
}

//关闭actor，异步函数等待actor主线程销毁。goroutine safe
func  (ac *Actor) Close()  {
	nlog.Debug("request Close:",ac.name)
	ac.cancel()
}

//actor结束，释放资源
func (ac *Actor) onDestroy()  {
	nlog.Debug("onDestroy:",ac.name)
	ac.subjects.Range(func(key, v interface{}) bool {
		v.(*nats.Subscription).Unsubscribe()
		return true
	})
	if ac.proc!=nil {
		ac.proc.OnDestroy()
	}
}