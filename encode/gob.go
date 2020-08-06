package encode

import (
	"fmt"
	"reflect"
	"sync"
)

var name2typeMap sync.Map
//注册类型
func RegisterName(value interface{})  {
	t := reflect.TypeOf(value).Elem()
	name := t.String()
	name2typeMap.Store(name,t)
}
//通过名字new对象
func NewObjectByName(name string) (interface{},error) {
	o,ok := name2typeMap.Load(name)
	if !ok {
		return nil, fmt.Errorf("unknow type %q,need RegisterName", name)
	}
	t := o.(reflect.Type)
	msg := reflect.New(t).Interface()
	return msg,nil
}
