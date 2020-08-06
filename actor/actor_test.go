package actor

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestEncoder(t *testing.T)  {
	nc, _ := nats.Connect(nats.DefaultURL)
	c, _ := nats.NewEncodedConn(nc, nats.JSON_ENCODER)

	defer c.Close()

	type person struct {
		Name     string
		Address  string
		Age      int
	}
	me := &person{Name: "derek", Age: 22, Address: "140 New Montgomery Street, San Francisco, CA"}

	// Simple Async Subscriber
	c.Subscribe("foo", func(s *person) {
		t.Logf("Received a message: %v\n", s)
	})

	// Simple Publisher
	c.Publish("foo", me)

	time.Sleep(time.Second)
}


func TestActor(t *testing.T)  {

	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		t.Fatal(err)
	}
	ac := NewActor("test",nc,nil)
	ac.Start()
	defer  ac.Close()

	ac.Subscribe("test2")
	go ac.Run()

	NewProxy("test2",nc).Tell("hello2")
	time.Sleep(time.Second)
}


type actorproc struct {

}
func (p *actorproc) OnStart() {
	fmt.Println("##actorproc OnStart")
}
func (p *actorproc) Receive(ctx Context) {
	rmsg := ctx.GetRawMsg()
	msg := ctx.Message()
	fmt.Println("##actorproc Receive:",msg)
	fmt.Println("##actorproc Receive raw:",string(rmsg.Data))
	switch m:=msg.(type) {
	case string:
		ctx.RespondMessage(m+" world")
	}
}
func (p *actorproc) OnDestroy() {
	fmt.Println("##actorproc OnDestroy")
}

func TestActorProc(t *testing.T)   {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		t.Fatal(err)
	}
	proc := actorproc{}
	ac := NewActor("test",nc,&proc)
	ac.Start()
	defer  ac.Close()

	assert.Nil(t,ac.Subscribe("test2"))
	go ac.Run()

	assert.Nil(t,ac.Tell("hello"))

	time.Sleep(time.Second)
}

func TestUnsub(t *testing.T)  {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		t.Fatal(err)
	}
	ac := NewActor("test",nc,nil)
	ac.Start()
	defer  ac.Close()

	ac.Subscribe("test2")
	//go ac.Run()


	pr := NewProxy("test2",nc)
	pr.Tell("hello1")
	ctx := ac.Read()
	msg := ctx.Message()
	assert.Equal(t,msg.(string),"hello1")

	t.Log("unsub...")
	ac.Unsubscribe("test2")
	//should not recv hello2
	pr.Tell("hello2")
	time.Sleep(time.Second)
	assert.Equal(t,len(ac.ch),0)
}

func TestProxy(t *testing.T) {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		t.Fatal(err)
	}
	ac := NewActor("test",nc,nil)
	ac.Start()
	defer ac.Close()

	//go ac.Run()

	//tell
	p := NewProxy(ac.subjectName(),nc)
	if errp := p.Tell("hello");errp!=nil {
		t.Fatal(errp)
	}
	ctx := ac.Read()
	assert.Equal(t,ctx.Message().(string),"hello")
}

func TestProxyRequest(t *testing.T) {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		t.Fatal(err)
	}
	proc := actorproc{}
	ac := NewActor("test",nc,&proc)
	ac.Start()
	defer ac.Close()

	go ac.Run()

	//proxy do something
	p := NewProxy(ac.subjectName(),nc)
	//request
	if msg,errp := p.Request("hello",time.Second*3);errp!=nil {
		t.Fatal(errp)
	} else {
		t.Log("response:",msg)
	}

	//async request
	ch,errAq:= p.AsyncRequest("hi",time.Second*3)
	assert.Nil(t,errAq)

	rsp := <-ch
	assert.Nil(t,rsp.Err)
	t.Log("async response:",rsp.Msg)
}