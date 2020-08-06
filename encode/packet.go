package encode

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"

	"errors"
	"fmt"

	"math"
	"github.com/golang/protobuf/proto"
	"reflect"
)

var default_order binary.ByteOrder = binary.BigEndian

func SetByteOrder(order binary.ByteOrder)  {
	default_order = order
}

func GetByteOrder() binary.ByteOrder {
	return default_order
}


//数据包
type Packet struct {
	bytes.Buffer
	order binary.ByteOrder
}

func NewPacket(data []byte) Packet {
	p := Packet{order:GetByteOrder()}
	if data!=nil {
		p.Buffer.Write(data)
	}
	return p
}

func (pk *Packet) ReadInt8() (int8, error) {
	b, err := pk.ReadByte()
	return int8(b), err
}

func (pk *Packet) ReadUInt8() (uint8, error) {
	return pk.ReadByte()
}

func (pk *Packet) ReadInt16() (int16, error) {
	buf := make([]byte, 2)

	n, err := pk.Read(buf)
	if err != nil {
		return 0, err
	}

	if n != 2 {
		return 0, errors.New("Read buf error")
	}

	return int16(pk.order.Uint16(buf)), nil
}

func (pk *Packet) ReadInt32() (int32, error) {
	buf := make([]byte, 4)

	n, err := pk.Read(buf)
	if err != nil {
		return 0, err
	}

	if n != 4 {
		return 0, errors.New("Read buf error")
	}

	return int32(pk.order.Uint32(buf)), nil
}

func (pk *Packet) ReadInt64() (int64, error) {
	buf := make([]byte, 8)

	n, err := pk.Read(buf)
	if err != nil {
		return 0, err
	}

	if n != 8 {
		return 0, errors.New("Read buf error")
	}

	return int64(pk.order.Uint64(buf)), nil
}

func (pk *Packet) ReadUInt16() (uint16, error) {
	buf := make([]byte, 2)

	n, err := pk.Read(buf)
	if err != nil {
		return 0, err
	}

	if n != 2 {
		return 0, errors.New("Read buf error")
	}

	return pk.order.Uint16(buf), nil
}

func (pk *Packet) ReadUInt32() (uint32, error) {
	buf := make([]byte, 4)

	n, err := pk.Read(buf)
	if err != nil {
		return 0, err
	}

	if n != 4 {
		return 0, errors.New("Read buf error")
	}

	return pk.order.Uint32(buf), nil
}

func (pk *Packet) ReadUInt64() (uint64, error) {
	buf := make([]byte, 8)

	n, err := pk.Read(buf)
	if err != nil {
		return 0, err
	}

	if n != 8 {
		return 0, errors.New("Read buf error")
	}

	return pk.order.Uint64(buf), nil
}

func (pk *Packet) ReadBytes() ([]byte, error) {

	n, err := pk.ReadUInt32()
	if err != nil {
		return nil, err
	}

	if n < 0 {
		return nil, errors.New("format error")
	}

	buf := make([]byte, n)
	rn, err := pk.Read(buf)
	if err != nil || rn != int(n) {
		return nil, errors.New("Read buf error")
	}

	return buf, nil

}

//长度+字节流
func (pk *Packet) ReadLString() (string, error) {
	bs, err := pk.ReadBytes()
	if err != nil {
		return "", err
	}

	return string(bs), err
}

func (pk *Packet) ReadFloat32() (float32, error) {

	d, err := pk.ReadUInt32()
	if err != nil {
		return 0, err
	}

	return math.Float32frombits(d), nil
}

func (pk *Packet) ReadFloat64() (float64, error) {

	d, err := pk.ReadUInt64()
	if err != nil {
		return 0, err
	}

	return math.Float64frombits(d), nil
}



func (pk *Packet) WriteInt8(x int8) {
	_ = pk.WriteByte(byte(x))
}

func (pk *Packet) WriteUInt8(b uint8) {
	_ = pk.WriteByte(b)
}

func (pk *Packet) WriteInt16(b int16) {
	buf := make([]byte, 2)
	pk.order.PutUint16(buf, uint16(b))
	pk.Write(buf)
}

func (pk *Packet) WriteInt32(b int32) {
	buf := make([]byte, 4)
	pk.order.PutUint32(buf, uint32(b))
	pk.Write(buf)
}

func (pk *Packet) WriteInt64(b int64) {
	buf := make([]byte, 8)
	pk.order.PutUint64(buf, uint64(b))
	pk.Write(buf)
}

func (pk *Packet) WriteUInt16(b uint16) {
	buf := make([]byte, 2)
	pk.order.PutUint16(buf, b)
	pk.Write(buf)
}

func (pk *Packet) WriteUInt32(b uint32) {
	buf := make([]byte, 4)
	pk.order.PutUint32(buf, b)
	pk.Write(buf)
}

func (pk *Packet) WriteUInt64(b uint64) {
	buf := make([]byte, 8)
	pk.order.PutUint64(buf, b)
	pk.Write(buf)
}

func (pk *Packet) WriteBytes(bs []byte) {
	pk.WriteUInt32(uint32(len(bs)))
	pk.Write(bs)
}

//长度+字节流
func (pk *Packet) WriteLString(s string) {
	pk.WriteUInt32(uint32(len(s)))
	pk.Write([]byte(s))
}

func (pk *Packet) WriteFloat32(f float32) {
	pk.WriteUInt32(math.Float32bits(f))
}

func (pk *Packet) WriteFloat64(f float64) {
	pk.WriteUInt64(math.Float64bits(f))
}


/**
 *写pb对象
 *对象需要提前proto.RegisterMapType注册
 *
**/
func (pk *Packet) WritePBObject(arg interface{}) error{
	pbm, isPB := arg.(proto.Message)
	if !isPB {
		return errors.New(fmt.Sprintf("%v is not protobuff struct",reflect.TypeOf(arg)))
	}
	t := reflect.TypeOf(arg).Elem()
	name := t.String()

	data,err := proto.Marshal(pbm)
	if err != nil {
		return err
	}
	pk.WriteLString(name)
	pk.WriteBytes(data)

	return nil
}

/**
 *读pb对象
 *对象需要提前proto.RegisterMapType注册
 *
**/
func (pk *Packet) ReadPBObject() (interface{}, error) {
	aname,errRead:= pk.ReadLString()
	if errRead != nil {
		return nil,errRead
	}
	//new struct
	t := proto.MessageType(aname)
	if t == nil {
		return nil, fmt.Errorf("any: message type %q isn't linked in", aname)
	}
	msg := reflect.New(t.Elem()).Interface().(proto.Message)
	//unmarshal
	data,errbt:= pk.ReadBytes()
	if errbt != nil {
		return nil,errbt
	}
	e := proto.Unmarshal(data,msg)
	return msg,e
}


/**
 *写gob对象
 *gob对象需要提前RegisterName注册
 *仅支持go
**/
func (pk *Packet) WriteGobObject(arg interface{}) error{
	t := reflect.TypeOf(arg).Elem()
	name := t.String()
	pk.WriteLString(name)
	en := gob.NewEncoder(&pk.Buffer)
	err := en.Encode(arg)
	if err != nil {
		return err
	}

	return nil
}



/**
 *读gob对象
 *gob对象需要提前RegisterName注册
 *仅支持go
**/
func (pk *Packet) ReadGobObject() (interface{}, error) {
	aname,errRead:= pk.ReadLString()
	if errRead != nil {
		return nil,errRead
	}
	//new struct
	msg,err := NewObjectByName(aname)
	if err != nil {
		return nil, err
	}
	//unmarshal
	dec := gob.NewDecoder(&pk.Buffer)
	e := dec.Decode(msg)
	return msg,e
}


