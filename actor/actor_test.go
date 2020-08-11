package actor

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
	. "github.com/smartystreets/goconvey/convey"
)

func TestEncoder(t *testing.T)  {
	Convey("TestEncoder", t,func() {
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
	})

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

	pr := NewProxy("test2",nc)
	pr.Tell("hello2")
	time.Sleep(time.Second)
}


type actorproc struct {

}
func (p *actorproc) OnStart() {
	fmt.Println("##actorproc OnStart")
}
func (p *actorproc) Receive(ctx Context) {
	msg := ctx.Message()
	fmt.Println("##actorproc Receive:",msg)
	switch m:=msg.(type) {
	case string:
		ctx.RespondMessage(m+" world")
	case *WatchTerminated:
		fmt.Println("recv Terminated:",m)
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

func TestKillActor(t *testing.T)  {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		t.Fatal(err)
	}
	proc := actorproc{}
	ac := NewActor("test",nc,&proc)
	ac.Start()
	defer ac.Close()

	go ac.Run()

	ac.Tell(&Kill{Reason:"you are dead",Who:ac.name})

	time.Sleep(time.Second)
}

func TestWatch(t *testing.T)  {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		t.Fatal(err)
	}

	//new ac1
	proc := actorproc{}
	ac := NewActor("test",nc,&proc)
	ac.Start()
	defer ac.Close()
	go ac.Run()

	//new ac2
	proc2 := actorproc{}
	ac2 := NewActor("test2",nc,&proc2)
	ac2.Start()
	defer ac2.Close()
	//go ac.Run()
	//watch
	ac2.Tell("hi")
	hi := ac2.Read()
	assert.Equal(t,hi.Message().(string),"hi")
	hi.Watch(ac.name)



	//new ac3
	proc3 := actorproc{}
	ac3:= NewActor("test3",nc,&proc3)
	ac3.Start()
	defer ac3.Close()
	//go ac.Run()
	//watch
	ac3.Tell("hi")
	hi3 := ac3.Read()
	assert.Equal(t,hi3.Message().(string),"hi")
	hi3.Watch(ac.name)
	time.Sleep(time.Second)
	hi3.Unwatch(ac.name)

	time.Sleep(time.Second)
	//kill
	ac.Tell(&Kill{Reason:"you are dead",Who:ac.name})
	//recv term
	term := ac2.Read()
	assert.Equal(t,term.Message().(*WatchTerminated).Who,"test")

	//ac3 should not recv WatchTerminated
	assert.Equal(t,len(ac3.ch),0)
}