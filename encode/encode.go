package encode

import (
	"fmt"
	"github.com/golang/protobuf/proto"
)

/**
 *
 *消息包结构：
 * 基础类型：类型id(1)  二进制数据(1,2,4,8)
 * pb类型: 类型id(1)  结构名(n) 二进制数据
 * gob类型：类型id(1)  结构名(n) 二进制数据
**/
//封包
func Encode(arg interface{}) ([]byte, error) {
	enc := NewPacket(nil)
	var err error
	switch arg.(type) {
	case bool:
		enc.WriteInt8(Bool)
		if arg.(bool) {
			enc.WriteInt8(1)
		} else {
			enc.WriteInt8(0)
		}
	case int8:
		enc.WriteInt8(Int8)
		enc.WriteInt8(arg.(int8))
	case uint8:
		enc.WriteInt8(UInt8)
		enc.WriteUInt8(arg.(uint8))
	case int16:
		enc.WriteInt8(Int16)
		enc.WriteInt16(arg.(int16))
	case uint16:
		enc.WriteInt8(UInt16)
		enc.WriteUInt16(arg.(uint16))
	case int32:
		enc.WriteInt8(Int32)
		enc.WriteInt32(arg.(int32))
	case uint32:
		enc.WriteInt8(UInt32)
		enc.WriteUInt32(arg.(uint32))
	case int64:
		enc.WriteInt8(Int64)
		enc.WriteInt64(arg.(int64))
	case uint64:
		enc.WriteInt8(UInt64)
		enc.WriteUInt64(arg.(uint64))
	case float32:
		enc.WriteInt8(Float32)
		enc.WriteFloat32(arg.(float32))
	case float64:
		enc.WriteInt8(Float64)
		enc.WriteFloat64(arg.(float64))

	case string:
		enc.WriteInt8(String)
		enc.WriteLString(arg.(string))

	case []byte:
		enc.WriteInt8(Bytes)
		enc.WriteBytes(arg.([]byte))
	default:
		pbm, isPB := arg.(proto.Message)
		if isPB {
			enc.WriteInt8(Proto)
			err = enc.WritePBObject(pbm)
		} else {
			enc.WriteInt8(Gob)
			err = enc.WriteGobObject(pbm)
		}
	}
	return enc.Bytes(),err
}

//解包
func Decode(data []byte) (interface{}, error)  {
	dec := NewPacket(data)
	var typ int8
	var err error
	var v interface{}

	if typ, err = dec.ReadInt8(); err != nil {
		return nil,err
	}

	switch typ {
	case Bool:
		v, err = dec.ReadInt8()
		if v.(int8) == 1 {
			v = true
		} else {
			v = false
		}
	case Int8:
		v, err = dec.ReadInt8()
	case UInt8:
		v, err = dec.ReadUInt8()
	case Int16:
		v, err = dec.ReadInt16()
	case UInt16:
		v, err = dec.ReadUInt16()
	case Int32:
		v, err = dec.ReadInt32()
	case UInt32:
		v, err = dec.ReadUInt32()
	case Int64:
		v, err = dec.ReadInt64()
	case UInt64:
		v, err = dec.ReadUInt64()
	case Float32:
		v, err = dec.ReadFloat32()
	case Float64:
		v, err = dec.ReadFloat64()
	case String:
		v, err = dec.ReadLString()
	case Bytes:
		v, err = dec.ReadBytes()
	case Proto:
		v, err = dec.ReadPBObject()
	case Gob:
		v, err = dec.ReadGobObject()
	default:
		return nil, fmt.Errorf("no support type:%v",typ)
	}
	return v,err
}