package main

import (
	"fmt"
	"os"
)

//meta信息的载体，代码生成的主体


//字段类型
const(
	eTypeInt8 = iota
	eTypeUint8
	eTypeInt16
	eTypeUint16
	eTypeInt32
	eTypeUint32
	eTypeInt64
	eTypeUint64
	eTypeFloat32
	eTypeFloat64
	eTypeEnum

	eTypeString
	eTypeStruct
	eTypeInt8Array
	eTypeUint8Array
	eTypeInt16Array
	eTypeUint16Array
	eTypeInt32Array
	eTypeUint32Array
	eTypeInt64Array
	eTypeUint64Array
	eTypeFloat32Array
	eTypeFloat64Array
	eTypeStringArray
	eTypeStructArray
	eTypeEnumArray
)

//枚举
type ProtoEnum struct{
	Name string
	Values map[string]int

}

//结构体字段
type ProtoField struct{
	Name string //字段名
	Index int //编码
	Type int //正则匹配到的类型枚举
	//TypeName string //类型，如vector<int>
	BaseTypeName string //基本类型，如int
	Dscr string //描述字段，需要进一步解析
}

//结构体
type ProtoStruct struct{
	Name string
	Values map[int]ProtoField
}

type CodeWriter struct{
	Deep int
	writer *os.File
}

func (self* CodeWriter) DeepIn(){
	self.Deep++
}

func(self* CodeWriter) DeepOut(){
	self.Deep--
	if self.Deep<0{
		self.Deep=0
	}
}

func (self *CodeWriter) WriteLine(a ...interface{}){
	s := ""
	for i:=0;i<self.Deep;i++ {
		s += "\t"
	}
	self.writer.WriteString(s + fmt.Sprintln(a...))
}

func NewCodeWriter(f *os.File) *CodeWriter{
	return &CodeWriter{0,f}
}