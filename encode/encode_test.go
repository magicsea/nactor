package encode

import (
	"github.com/magicsea/nactor/pb"
	. "github.com/smartystreets/goconvey/convey"
	"reflect"
	"testing"
)

type St struct {
	A int
	S string
}

func TestPacket(t *testing.T)  {
	Convey("TestPacket",t, func() {
		Convey("readpb",func() {
			pk := NewPacket(nil)
			msg := pb.TellRequest{"hello",3}
			errEn := pk.WritePBObject(&msg)
			So(errEn,ShouldBeNil)
			pkr := NewPacket(pk.Bytes())

			msg2,errDe:= pkr.ReadPBObject()
			So(errDe,ShouldBeNil)
			So(reflect.DeepEqual(&msg,msg2),ShouldBeTrue)
		})

		Convey("readgob",func() {
			RegisterName((*St)(nil))
			pk := NewPacket(nil)
			msg := St{1,"kkk"}
			errEn := pk.WriteGobObject(&msg)
			So(errEn,ShouldBeNil)
			pkr := NewPacket(nil)
			pkr.Write(pk.Bytes())

			msg2,errDe:= pkr.ReadGobObject()
			So(errDe,ShouldBeNil)
			So(reflect.DeepEqual(&msg,msg2),ShouldBeTrue)
		})
	})
}


func TestGob(t *testing.T)  {

	RegisterName((*St)(nil))
	Convey("TestPacket",t, func() {
		o,errEn := NewObjectByName("encode.St")
		So(errEn,ShouldBeNil)
		t.Logf("%+v",o)
	})
}
