package main

import (
	"regexp"
	"errors"
	"strconv"
	"strings"
	"os"
	"path"
	"time"
	"sort"
)

//文件解析状态
const(
	none = iota
	beginEnum
	beginStruct
)

type Lexer interface{
	ReadLine(string)error
	GenCppCode(h,cpp string)
	GenCsharpCode(cs,s string)
}


type lexer struct{
	myState int
	myKey string

	structlst []string
	structs map[string]ProtoStruct
	enums map[string]ProtoEnum

	regEnum *regexp.Regexp
	regEnumField *regexp.Regexp
	regStruct *regexp.Regexp
	regStructField *regexp.Regexp
	regSpace *regexp.Regexp
	regComment *regexp.Regexp
	regEnd *regexp.Regexp
}

func (self* lexer) Init(){
	self.regEnum = regexp.MustCompile(`\s*enum\s+(\w+)\s*\{\s*$`)
	self.regStruct = regexp.MustCompile(`\s*struct\s+(\w+)\s*\{\s*$`)
	self.regComment = regexp.MustCompile(`^\s*//.*$`)
	self.regSpace = regexp.MustCompile(`^\s*$`)
	self.regStructField = regexp.MustCompile(`^\s*(\w+)\s+(\w+\+?)\s*=\s*(\d+)\s*("(.*)")?\s*$`)
	self.regEnumField = regexp.MustCompile(`^\s*(\w+)\s*=\s*(\d+)\s*$`)
	self.regEnd = regexp.MustCompile(`^\s*}\s*$`)

	self.enums = make(map[string]ProtoEnum)
	self.structs = make(map[string]ProtoStruct)
	self.structlst = make([]string,0)
}

func (self *lexer) ReadLine( s string) error{
	if self.regSpace.MatchString(s){
		return nil
	}

	if self.regComment.MatchString(s){
		return nil
	}

	if self.myState == none {
		if self.regEnum.MatchString(s){
			self.myState = beginEnum
			matches := self.regEnum.FindStringSubmatch(s)
			self.myKey = matches[1]
			self.enums[self.myKey]= ProtoEnum{Name:self.myKey, Values:make(map[string]int)}
			return nil
		} else 	if self.regStruct.MatchString(s){
			self.myState = beginStruct
			matches := self.regStruct.FindStringSubmatch(s)
			self.myKey = matches[1]

			self.structlst = append(self.structlst,self.myKey)

			self.structs[self.myKey] = ProtoStruct{Name:self.myKey, Values:make(map[int]ProtoField)}
			return nil
		} else {
			return errors.New("语法错误："+s)
		}
	}


	if self.myState == beginEnum{
		if self.regEnd.MatchString(s){
			self.myState=none
			return nil
		}
		matches := self.regEnumField.FindStringSubmatch(s)
		if matches == nil{
			return errors.New("枚举字段错误："+s)
		}
		v,err := strconv.Atoi(matches[2])
		if err!=nil{
			return err
		}
		self.enums[self.myKey].Values[matches[1]]=v
		return nil
	}

	if self.myState == beginStruct{
		if self.regEnd.MatchString(s){
			self.myState=none
			return nil
		}

		matches := self.regStructField.FindStringSubmatch(s)
		if matches == nil{
			return errors.New("结构体字段错误："+s)
		}

		i,err := strconv.Atoi(matches[3])
		if err != nil{
			return err
		}

		t,tname := self.GetFieldType(matches[2])
		if t == -1{
			return errors.New("类型错误："+matches[2])
		}

		if _,ok:=self.structs[self.myKey].Values[i];ok{
			return errors.New("结构体序号重复："+s)
		}


		self.structs[self.myKey].Values[i]=ProtoField{
			Name:matches[1],
			Type:t,
			Index:i,
			BaseTypeName:tname,
			Dscr:matches[4]}
		return nil
	}

	return errors.New("无法解析的行："+s)
}

func (self *lexer)GetFieldType(s string) (int, string){

	switch s{
	case "int8":
		return eTypeInt8,s
	case "uint8":
		return eTypeUint8,s
	case "int16":
		return eTypeInt16,s
	case "uint16":
		return eTypeUint16,s
	case "int32":
		return eTypeInt32,s
	case "uint32":
		return eTypeUint32,s
	case "float32":
		return eTypeFloat32,s
	case "float64":
		return eTypeFloat64,s
	case "int64":
		return eTypeInt64,s
	case "uint64":
		return eTypeUint64,s
	case "int":
		return eTypeInt32,s
	case "string":
		return eTypeString,s
	case "float":
		return eTypeFloat32,s
	case "int+":
		return eTypeInt32Array,"int"
	case "float+":
		return eTypeFloat32Array,"float"
	case "string+":
		return  eTypeStringArray,"string"
	case "int8+":
		return eTypeInt8,"int8_t"
	case "uint8+":
		return eTypeUint8,"uint8_t"
	case "int16+":
		return eTypeInt16,"int16_t"
	case "uint16+":
		return eTypeUint16,"uint16_t"
	case "int32+":
		return eTypeInt32,"int32_t"
	case "uint32+":
		return eTypeUint32,"uint32_t"
	case "float32+":
		return eTypeFloat32,"float"
	case "float64+":
		return eTypeFloat64,"double"
	case "int64+":
		return eTypeInt64,"int64_t"
	case "uint64+":
		return eTypeUint64Array,"uint64_t"

	}

	overlap := false
	if strings.HasSuffix(s,"+"){
		overlap=true
		s = strings.TrimRight(s,"+")
	}
	if _,ok := self.enums[s]; ok {
		if overlap{
			return eTypeEnumArray,s
		}
		return eTypeEnum,s
	}

	if _,ok := self.structs[s]; ok{
		if overlap{
			return eTypeStructArray,s
		}
		return eTypeStruct,s
	}


	return -1,s
}

func (self *lexer) GetCppTypeName(t int, tname string) string{
	switch t {
	case eTypeInt8:
		return "int8_t"
	case eTypeUint8:
		return "uint8_t"
	case eTypeInt16:
		return "int16_t"
	case eTypeUint16:
		return "uint16_t"
	case eTypeInt32:
		return "int32_t"
	case eTypeUint32:
		return "uint32_t"
	case eTypeInt64:
		return "int64_t"
	case eTypeUint64:
		return "uint64_t"
	case eTypeFloat32:
		return "float"
	case eTypeFloat64:
		return "double"
	case eTypeString:
		return "string"
	case eTypeInt8Array:
		return "vector<int8_t>"
	case eTypeUint8Array:
		return "vector<uint8_t>"
	case eTypeInt16Array:
		return "vector<int16_t>"
	case eTypeUint16Array:
		return "vector<uint16_t>"
	case eTypeInt32Array:
		return "vector<int32_t>"
	case eTypeUint32Array:
		return "vector<uint32_t>"
	case eTypeFloat32Array:
		return "vector<float>"
	case eTypeFloat64Array:
		return "vector<double>"
	case eTypeStringArray:
		return "vector<string>"
	case eTypeStruct:
		if s,ok := self.structs[tname];ok{
			return s.Name
		}
		return tname
	case eTypeStructArray:
		if s,ok := self.structs[tname];ok{
			return "vector<"+s.Name+">"
		}
		return tname
	case eTypeEnum:
		if s,ok := self.enums[tname];ok{
			return s.Name
		}
		return tname
	case eTypeEnumArray:
		if s,ok:=self.enums[tname];ok{
			return "vector<"+s.Name+">"
		}
		return tname
	}

	return tname
}

func (self *lexer) GetCsTypeName(t int, tname string) string{
	switch t {
	case eTypeInt8:
		return "sbyte"
	case eTypeUint8:
		return "byte"
	case eTypeInt16:
		return "short"
	case eTypeUint16:
		return "ushort"
	case eTypeInt32:
		return "int"
	case eTypeUint32:
		return "uint"
	case eTypeInt64:
		return "long"
	case eTypeUint64:
		return "ulong"
	case eTypeFloat32:
		return "float"
	case eTypeFloat64:
		return "double"
	case eTypeString:
		return "string"
	case eTypeInt8Array:
		return "List<sbyte>"
	case eTypeUint8Array:
		return "List<byte>"
	case eTypeInt16Array:
		return "List<short>"
	case eTypeUint16Array:
		return "List<ushort>"
	case eTypeInt32Array:
		return "List<int>"
	case eTypeUint32Array:
		return "List<uint>"
	case eTypeFloat32Array:
		return "List<float>"
	case eTypeFloat64Array:
		return "List<double>"
	case eTypeStringArray:
		return "List<string>"
	case eTypeStruct:
		if s,ok := self.structs[tname];ok{
			return s.Name
		}
		return tname
	case eTypeStructArray:
		if s,ok := self.structs[tname];ok{
			return "List<"+s.Name+">"
		}
		return tname
	case eTypeEnum:
		if s,ok := self.enums[tname];ok{
			return s.Name
		}
		return tname
	case eTypeEnumArray:
		if s,ok:=self.enums[tname];ok{
			return "List<"+s.Name+">"
		}
		return tname
	}

	return tname
}

func (self* lexer)GenHeadFile(h string){
	hfile ,err := os.Create(h)
	if err != nil{
		panic(err)
	}

	defer hfile.Close()

	cw := NewCodeWriter(hfile)
	cw.WriteLine("//THIS FILE IS GENERATED BY STRUCTPROTO, DO NOT EDIT IT!!!")
	cw.WriteLine("//GENERATE TIME : ", time.Now())
	cw.WriteLine("#pragma once\n")
	cw.WriteLine("#include <string>")
	cw.WriteLine("#include <vector>")
	cw.WriteLine("using namespace std;\n")
	//枚举
	for k,v := range self.enums{
		cw.WriteLine("enum "+k+" {")
		cw.DeepIn()
		for vk,vv := range v.Values{

			cw.WriteLine(vk," = ",vv,",")
		}
		cw.DeepOut()
		cw.WriteLine("};\n")
	}

	//结构体
	for _,key := range self.structlst{
		if v,ok := self.structs[key];ok{
			cw.WriteLine("struct "+key+" {")
			cw.DeepIn()

			keylst := make([]int,0)
			for k := range v.Values{
				keylst=append(keylst,k)
			}
			sort.Ints(keylst)
			for _,k := range keylst{
				vv := v.Values[k]
				cw.WriteLine(self.GetCppTypeName(vv.Type, vv.BaseTypeName)," ",vv.Name,";")

			}
			cw.WriteLine("int serialize(void* buff, int len) const;")
			cw.WriteLine("void deserialize(void* buff, int len);")
			cw.DeepOut()
			cw.WriteLine("};\n")
		}
	}



}

func (self* lexer) GenCppFile(cpp string){
	cppfile, err := os.Create(cpp)
	if err != nil {
		panic(err)
	}
	defer cppfile.Close()

	cw := NewCodeWriter(cppfile)

	_,name := path.Split(cpp)
	hname := strings.TrimSuffix(name,".cpp")+".h"
	cw.WriteLine("#include \""+hname+"\"")
	cw.WriteLine("#include \"BuffBuilder.h\"")
	cw.WriteLine("#include \"BuffParser.h\"")
	cw.WriteLine("\n")

	for k,v := range self.structs{
		cw.WriteLine("int "+k+"::serialize(void* buff, int len) const")
		cw.WriteLine("{")
		cw.DeepIn()
		cw.WriteLine("BuffBuilder bb(buff, len);")
		cw.WriteLine("uint16_t _member_count = ",len(v.Values),";")
		cw.WriteLine("bb.Push(_member_count);")
		cw.WriteLine("uint16_t _member_code = 0;")
		cw.WriteLine("uint8_t _member_type=0;")

		for _,vv := range v.Values{
			cw.WriteLine("_member_code = ",vv.Index,";")
			cw.WriteLine("bb.Push(_member_code);")
			cw.WriteLine("_member_type = ",vv.Type,";")
			cw.WriteLine("bb.Push(_member_type);")
			switch vv.Type{
			case eTypeInt8:
				cw.WriteLine("bb.PushInt8("+vv.Name+");")
			case eTypeUint8:
				cw.WriteLine("bb.PushUInt8("+vv.Name+");")
			case eTypeInt16:
				cw.WriteLine("bb.PushInt16("+vv.Name+");")
			case eTypeUint16:
				cw.WriteLine("bb.PushUInt16("+vv.Name+");")
			case eTypeEnum:
				cw.WriteLine("bb.PushInt32((int)"+vv.Name+");")
			case eTypeInt32:
				cw.WriteLine("bb.PushInt32("+vv.Name+");")
			case eTypeUint32:
				cw.WriteLine("bb.PushUInt32("+vv.Name+");")
			case eTypeInt64:
				cw.WriteLine("bb.PushInt64("+vv.Name+");")
			case eTypeUint64:
				cw.WriteLine("bb.PushUInt64("+vv.Name+");")
			case eTypeFloat32:
				cw.WriteLine("bb.PushFloat32("+vv.Name+");")
			case eTypeFloat64:
				cw.WriteLine("bb.PushFloat64("+vv.Name+");")
			case eTypeString:
				cw.WriteLine("bb.PushString("+vv.Name+");")
			case eTypeStruct:
				cw.WriteLine("uint16_t& elemlen = *(uint16_t*)bb.Cursor();")
				cw.WriteLine("elemlen = 0;")
				cw.WriteLine("bb.Push(elemlen);")
				cw.WriteLine("elemlen = ",vv.Name,".serialize(bb.Cursor(), bb.SizeToPush());")
				cw.WriteLine("bb.SkipSize(elemlen);")
			case eTypeStructArray:
				cw.WriteLine("bb.PushUInt16((uint16_t)",vv.Name,".size());")
				cw.WriteLine("for(int n=0; n < ",vv.Name,".size(); n++){")
				cw.DeepIn()
				cw.WriteLine("uint16_t& elemlen = *(uint16_t*)bb.Cursor();")
				cw.WriteLine("elemlen = 0;")
				cw.WriteLine("bb.Push(elemlen);")
				cw.WriteLine("elemlen = ",vv.Name,"[n].serialize(bb.Cursor(), bb.SizeToPush());")
				cw.WriteLine("bb.SkipSize(elemlen);")
				cw.DeepOut()
				cw.WriteLine("}\n")
			case eTypeInt8Array:
				cw.WriteLine("bb.PushUInt16((uint16_t)",vv.Name,".size());")
				cw.WriteLine("for(int n=0; n < ",vv.Name,".size(); n++){")
				cw.DeepIn()
				cw.WriteLine("bb.PushInt8(",vv.Name,"[n]);")
				cw.DeepOut()
				cw.WriteLine("}\n")
			case eTypeUint8Array:
				cw.WriteLine("bb.PushUInt16((uint16_t)",vv.Name,".size());")
				cw.WriteLine("for(int n=0; n < ",vv.Name,".size(); n++){")
				cw.DeepIn()
				cw.WriteLine("bb.PushUInt8(",vv.Name,"[n]);")
				cw.DeepOut()
				cw.WriteLine("}\n")
			case eTypeInt16Array:
				cw.WriteLine("bb.PushUInt16((uint16_t)",vv.Name,".size());")
				cw.WriteLine("for(int n=0; n < ",vv.Name,".size(); n++){")
				cw.DeepIn()
				cw.WriteLine("bb.PushInt16(",vv.Name,"[n]);")
				cw.DeepOut()
				cw.WriteLine("}\n")
			case eTypeUint16Array:
				cw.WriteLine("bb.PushUInt16((uint16_t)",vv.Name,".size());")
				cw.WriteLine("for(int n=0; n < ",vv.Name,".size(); n++){")
				cw.DeepIn()
				cw.WriteLine("bb.PushUInt16(",vv.Name,"[n]);")
				cw.DeepOut()
				cw.WriteLine("}\n")
			case eTypeEnumArray:
				cw.WriteLine("bb.PushUInt16((uint16_t)",vv.Name,".size());")
				cw.WriteLine("for(int n=0; n < ",vv.Name,".size(); n++){")
				cw.DeepIn()
				cw.WriteLine("bb.PushInt32((int)",vv.Name,"[n]);")
				cw.DeepOut()
				cw.WriteLine("}\n")
			case eTypeInt32Array:
				cw.WriteLine("bb.PushUInt16((uint16_t)",vv.Name,".size());")
				cw.WriteLine("for(int n=0; n < ",vv.Name,".size(); n++){")
				cw.DeepIn()
				cw.WriteLine("bb.PushInt32(",vv.Name,"[n]);")
				cw.DeepOut()
				cw.WriteLine("}\n")
			case eTypeUint32Array:
				cw.WriteLine("bb.PushUInt16((uint16_t)",vv.Name,".size());")
				cw.WriteLine("for(int n=0; n < ",vv.Name,".size(); n++){")
				cw.DeepIn()
				cw.WriteLine("bb.PushUInt32(",vv.Name,"[n]);")
				cw.DeepOut()
				cw.WriteLine("}\n")
			case eTypeInt64Array:
				cw.WriteLine("bb.PushUInt16((uint16_t)",vv.Name,".size());")
				cw.WriteLine("for(int n=0; n < ",vv.Name,".size(); n++){")
				cw.DeepIn()
				cw.WriteLine("bb.PushInt64(",vv.Name,"[n]);")
				cw.DeepOut()
				cw.WriteLine("}\n")
			case eTypeUint64Array:
				cw.WriteLine("bb.PushUInt16((uint16_t)",vv.Name,".size());")
				cw.WriteLine("for(int n=0; n < ",vv.Name,".size(); n++){")
				cw.DeepIn()
				cw.WriteLine("bb.PushUInt64(",vv.Name,"[n]);")
				cw.DeepOut()
				cw.WriteLine("}\n")
			case eTypeFloat32Array:
				cw.WriteLine("bb.PushUInt16((uint16_t)",vv.Name,".size());")
				cw.WriteLine("for(int n=0; n < ",vv.Name,".size(); n++){")
				cw.DeepIn()
				cw.WriteLine("bb.PushFloat32(",vv.Name,"[n]);")
				cw.DeepOut()
				cw.WriteLine("}\n")
			case eTypeFloat64Array:
				cw.WriteLine("bb.PushUInt16((uint16_t)",vv.Name,".size());")
				cw.WriteLine("for(int n=0; n < ",vv.Name,".size(); n++){")
				cw.DeepIn()
				cw.WriteLine("bb.PushFloat64(",vv.Name,"[n]);")
				cw.DeepOut()
				cw.WriteLine("}\n")
			case eTypeStringArray:
				cw.WriteLine("bb.PushUInt16((uint16_t)",vv.Name,".size());")
				cw.WriteLine("for(int n=0; n < ",vv.Name,".size(); n++){")
				cw.DeepIn()
				cw.WriteLine("bb.PushString(",vv.Name,"[n]);")
				cw.DeepOut()
				cw.WriteLine("}\n")

			}
		}

		cw.WriteLine("return bb.Size();")
		cw.DeepOut()
		cw.WriteLine("}")

		//deserailize
		cw.WriteLine("void "+k+"::deserialize(void* buff, int len)")
		cw.WriteLine("{")
		cw.DeepIn()
		cw.WriteLine("BuffParser bp(buff, len);")
		cw.WriteLine("uint16_t _member_count = bp.PopUInt16();")
		cw.WriteLine("uint16_t _member_code = 0;")
		cw.WriteLine("uint8_t _member_type = 0;")
		cw.WriteLine("for(uint16_t i=0;i<_member_count;i++)")
		cw.WriteLine("{")
		cw.DeepIn()

		cw.WriteLine("_member_code = bp.PopUInt16();")
		cw.WriteLine("_member_type = bp.PopUInt8();")
		cw.WriteLine("switch(_member_code){")
		cw.DeepIn()
		for _,vv := range v.Values {

			cw.WriteLine("case ",vv.Index,":")
			switch vv.Type{
			case eTypeInt8:
				cw.WriteLine(vv.Name+" = bp.PopInt8();break;")
			case eTypeUint8:
				cw.WriteLine(vv.Name+"  = bp.PopUInt8();break;")
			case eTypeInt16:
				cw.WriteLine(vv.Name+" = bp.PopInt16();break;")
			case eTypeUint16:
				cw.WriteLine(vv.Name+" = bp.PopUInt16();break;")
			case eTypeInt32:
				cw.WriteLine(vv.Name+" = bp.PopInt32();break;")
			case eTypeUint32:
				cw.WriteLine(vv.Name+" = bp.PopUInt32();break;")
			case eTypeFloat32:
				cw.WriteLine(vv.Name+" = bp.Float32();break;")
			case eTypeFloat64:
				cw.WriteLine(vv.Name+" = bp.Float64();break;")
			case eTypeString:
				cw.WriteLine(vv.Name+" = bp.PopString();break;")
			case eTypeEnum:
				cw.WriteLine(vv.Name+" = (",vv.BaseTypeName,")bp.PopInt32();break;")//TODO 检查值合法
			case eTypeInt8Array:
				cw.WriteLine("{")
				cw.DeepIn()
				cw.WriteLine("uint16_t arraylen = bp.PopUInt16();")
				cw.WriteLine("for (;arraylen>0;arraylen--){")
				cw.DeepIn()
				cw.WriteLine(vv.BaseTypeName+" element = bp.PopInt8();")
				cw.WriteLine(vv.Name+".push_back(element);")
				cw.DeepOut()
				cw.WriteLine("}")
				cw.DeepOut()
				cw.WriteLine("}")
				cw.WriteLine("break;")
			case eTypeUint8Array:
				cw.WriteLine("{")
				cw.DeepIn()
				cw.WriteLine("uint16_t arraylen = bp.PopUInt16();")
				cw.WriteLine("for (;arraylen>0;arraylen--){")
				cw.DeepIn()
				cw.WriteLine(vv.BaseTypeName+" element = bp.PopUInt8();")
				cw.WriteLine(vv.Name+".push_back(element);")
				cw.DeepOut()
				cw.WriteLine("}")
				cw.DeepOut()
				cw.WriteLine("}")
				cw.WriteLine("break;")
			case eTypeInt16Array:
				cw.WriteLine("{")
				cw.DeepIn()
				cw.WriteLine("uint16_t arraylen = bp.PopUInt16();")
				cw.WriteLine("for (;arraylen>0;arraylen--){")
				cw.DeepIn()
				cw.WriteLine(vv.BaseTypeName+" element = bp.PopInt16();")
				cw.WriteLine(vv.Name+".push_back(element);")
				cw.DeepOut()
				cw.WriteLine("}")
				cw.DeepOut()
				cw.WriteLine("}")
				cw.WriteLine("break;")
			case eTypeUint16Array:
				cw.WriteLine("{")
				cw.DeepIn()
				cw.WriteLine("uint16_t arraylen = bp.PopUInt16();")
				cw.WriteLine("for (;arraylen>0;arraylen--){")
				cw.DeepIn()
				cw.WriteLine(vv.BaseTypeName+" element = bp.PopUInt16();")
				cw.WriteLine(vv.Name+".push_back(element);")
				cw.DeepOut()
				cw.WriteLine("}")
				cw.DeepOut()
				cw.WriteLine("}")
				cw.WriteLine("break;")
			case eTypeInt32Array:
				cw.WriteLine("{")
				cw.DeepIn()
				cw.WriteLine("uint16_t arraylen = bp.PopUInt16();")
				cw.WriteLine("for (;arraylen>0;arraylen--){")
				cw.DeepIn()
				cw.WriteLine(vv.BaseTypeName+" element = bp.PopInt32();")
				cw.WriteLine(vv.Name+".push_back(element);")
				cw.DeepOut()
				cw.WriteLine("}")
				cw.DeepOut()
				cw.WriteLine("}")
				cw.WriteLine("break;")
			case eTypeUint32Array:
				cw.WriteLine("{")
				cw.DeepIn()
				cw.WriteLine("uint16_t arraylen = bp.PopUInt16();")
				cw.WriteLine("for (;arraylen>0;arraylen--){")
				cw.DeepIn()
				cw.WriteLine(vv.BaseTypeName+" element = bp.PopUInt32();")
				cw.WriteLine(vv.Name+".push_back(element);")
				cw.DeepOut()
				cw.WriteLine("}")
				cw.DeepOut()
				cw.WriteLine("}")
				cw.WriteLine("break;")
			case eTypeFloat32Array:
				cw.WriteLine("{")
				cw.DeepIn()
				cw.WriteLine("uint16_t arraylen = bp.PopUInt16();")
				cw.WriteLine("for (;arraylen>0;arraylen--){")
				cw.DeepIn()
				cw.WriteLine(vv.BaseTypeName+" element = bp.PopFloat32();")
				cw.WriteLine(vv.Name+".push_back(element);")
				cw.DeepOut()
				cw.WriteLine("}")
				cw.DeepOut()
				cw.WriteLine("}")
				cw.WriteLine("break;")
			case eTypeFloat64Array:
				cw.WriteLine("{")
				cw.DeepIn()
				cw.WriteLine("uint16_t arraylen = bp.PopUInt16();")
				cw.WriteLine("for (;arraylen>0;arraylen--){")
				cw.DeepIn()
				cw.WriteLine(vv.BaseTypeName+" element = bp.PopFloat64();")
				cw.WriteLine(vv.Name+".push_back(element);")
				cw.DeepOut()
				cw.WriteLine("}")
				cw.DeepOut()
				cw.WriteLine("}")
				cw.WriteLine("break;")
			case eTypeEnumArray:
				cw.WriteLine("{")
				cw.DeepIn()
				cw.WriteLine("uint16_t arraylen = bp.PopUInt16();")
				cw.WriteLine("for (;arraylen>0;arraylen--){")
				cw.DeepIn()
				cw.WriteLine(vv.BaseTypeName+" element = (",vv.BaseTypeName,")bp.PopInt32();")
				cw.WriteLine(vv.Name+".push_back(element);")
				cw.DeepOut()
				cw.WriteLine("}")
				cw.DeepOut()
				cw.WriteLine("}")
				cw.WriteLine("break;")
			case eTypeStruct:
				cw.WriteLine("{")
				cw.DeepIn()
				cw.WriteLine("uint16_t elemlen = bp.PopUInt16();")
				cw.WriteLine(vv.BaseTypeName+" element {};")
				cw.WriteLine("element.deserialize(bp.Cursor(), bp.SizeToPop());")
				cw.WriteLine("bp.SkipSize(elemlen);")
				cw.DeepOut()
				cw.WriteLine("}")
				cw.WriteLine("break;")
			case eTypeStructArray:
				cw.WriteLine("{")
				cw.DeepIn()
				cw.WriteLine("uint16_t arraylen = bp.PopUInt16();")
				cw.WriteLine("for (;arraylen>0;arraylen--){")
				cw.DeepIn()
				cw.WriteLine("uint16_t elemlen = bp.PopUInt16();")
				cw.WriteLine(vv.BaseTypeName+" element{};")
				cw.WriteLine("element.deserialize(bp.Cursor(), bp.SizeToPop());")
				cw.WriteLine(vv.Name+".push_back(element);")
				cw.WriteLine("bp.SkipSize(elemlen);")
				cw.DeepOut()
				cw.WriteLine("}")
				cw.DeepOut()
				cw.WriteLine("}")
				cw.WriteLine("break;")
			default:
				panic(errors.New("unknown type"))
			}

		}
		cw.WriteLine("default:")
		cw.WriteLine("bp.SkipType((eMemberType)_member_type);break;")
		cw.DeepOut()
		cw.WriteLine("}")//switch
		cw.DeepOut()
		cw.WriteLine("}")//for
		cw.DeepOut()
		cw.WriteLine("}\n")
	}

}

func (self *lexer)GenCppCode(h,cpp string){

	self.GenHeadFile(h)
	self.GenCppFile(cpp)
}

func (self *lexer)GenCsharpCode(cs,s string){
	csfile, err := os.Create(cs)
	if err != nil{
		panic(err)
	}

	defer csfile.Close()

	s = strings.Split(s,".")[0]

	cw := NewCodeWriter(csfile)
	cw.WriteLine("//THIS FILE IS GENERATED BY STRUCTPROTO, DO NOT EDIT IT!!!")
	cw.WriteLine("//GENERATE TIME : ", time.Now())
	cw.WriteLine("using System; ")
	cw.WriteLine("using System.Collections.Generic;")
	cw.WriteLine("using StructProtocol;")
	cw.WriteLine("\n")
	cw.WriteLine("namespace ",s,"{")
	cw.DeepIn()
	//枚举
	for k,v := range self.enums{
		cw.WriteLine("enum "+k+" {")
		cw.DeepIn()
		for vk,vv := range v.Values{

			cw.WriteLine(vk," = ",vv,",")
		}
		cw.DeepOut()
		cw.WriteLine("}\n")
	}

	//结构体
	for _,key := range self.structlst{
		if v,ok := self.structs[key];ok{
			cw.WriteLine("class "+key+": IStruct {")
			cw.DeepIn()

			keylst := make([]int,0)
			for k := range v.Values{
				keylst=append(keylst,k)
			}
			sort.Ints(keylst)
			for _,k := range keylst{
				vv := v.Values[k]
				if vv.Type < eTypeInt8Array {
					cw.WriteLine("public ", self.GetCsTypeName(vv.Type, vv.BaseTypeName), " ", vv.Name, ";")
				}else{
					cw.WriteLine("public ", self.GetCsTypeName(vv.Type,vv.BaseTypeName)," ",vv.Name,"= new ",self.GetCsTypeName(vv.Type,vv.BaseTypeName),"();")
				}
			}
			cw.WriteLine("public int serialize(byte[] buff) {")
			cw.DeepIn()
			cw.WriteLine("BuffBuilder bb = new BuffBuilder(buff);")
			cw.WriteLine("ushort _member_count = ",len(v.Values),";")
			cw.WriteLine("bb.PushUInt16(_member_count);")
			cw.WriteLine("ushort _member_code = 0;")
			cw.WriteLine("byte _member_type=0;")

			for _,vv := range v.Values{
				cw.WriteLine("_member_code = ",vv.Index,";")
				cw.WriteLine("bb.PushUInt16(_member_code);")
				cw.WriteLine("_member_type = ",vv.Type,";")
				cw.WriteLine("bb.PushUInt8(_member_type);")
				switch vv.Type{
				case eTypeInt8:
					cw.WriteLine("bb.PushInt8("+vv.Name+");")
				case eTypeUint8:
					cw.WriteLine("bb.PushUInt8("+vv.Name+");")
				case eTypeInt16:
					cw.WriteLine("bb.PushInt16("+vv.Name+");")
				case eTypeUint16:
					cw.WriteLine("bb.PushUInt16("+vv.Name+");")
				case eTypeEnum:
					cw.WriteLine("bb.PushInt32((int)"+vv.Name+");")
				case eTypeInt32:
					cw.WriteLine("bb.PushInt32("+vv.Name+");")
				case eTypeUint32:
					cw.WriteLine("bb.PushUint32("+vv.Name+");")
				case eTypeInt64:
					cw.WriteLine("bb.PushInt64("+vv.Name+");")
				case eTypeUint64:
					cw.WriteLine("bb.PushUint64("+vv.Name+");")
				case eTypeFloat32:
					cw.WriteLine("bb.PushFloat("+vv.Name+");")
				case eTypeFloat64:
					cw.WriteLine("bb.PushDouble("+vv.Name+");")
				case eTypeString:
					cw.WriteLine("bb.PushString("+vv.Name+");")
				case eTypeStruct:
					cw.WriteLine("byte[] _buf = new byte[4096];")
					cw.WriteLine("ushort elemlen = (ushort)",vv.Name,".serialize(_buf);")
					cw.WriteLine("bb.PushUInt16(elemlen);")
					cw.WriteLine("bb.PushRangeBuff(_buf,elemlen);")
				case eTypeStructArray:
					cw.WriteLine("bb.PushUInt16((ushort)",vv.Name,".Count);")
					cw.WriteLine("for(int n=0; n < ",vv.Name,".Count; n++){")
					cw.DeepIn()
					cw.WriteLine("byte[] _buf = new byte[4096];")
					cw.WriteLine("ushort elemlen = (ushort)",vv.Name,"[n].serialize(_buf);")
					cw.WriteLine("bb.PushUInt16(elemlen);")
					cw.WriteLine("bb.PushRangeBuff(_buf,elemlen);")
					cw.DeepOut()
					cw.WriteLine("}\n")
				case eTypeInt8Array:
					cw.WriteLine("bb.PushUInt16((ushort)",vv.Name,".Count);")
					cw.WriteLine("for(int n=0; n < ",vv.Name,".Count; n++){")
					cw.DeepIn()
					cw.WriteLine("bb.PushInt8(",vv.Name,"[n]);")
					cw.DeepOut()
					cw.WriteLine("}\n")
				case eTypeUint8Array:
					cw.WriteLine("bb.PushUInt16((ushort)",vv.Name,".Count);")
					cw.WriteLine("for(int n=0; n < ",vv.Name,".Count; n++){")
					cw.DeepIn()
					cw.WriteLine("bb.PushUInt8(",vv.Name,"[n]);")
					cw.DeepOut()
					cw.WriteLine("}\n")
				case eTypeInt16Array:
					cw.WriteLine("bb.PushUInt16((ushort)",vv.Name,".Count);")
					cw.WriteLine("for(int n=0; n < ",vv.Name,".Count; n++){")
					cw.DeepIn()
					cw.WriteLine("bb.PushInt16(",vv.Name,"[n]);")
					cw.DeepOut()
					cw.WriteLine("}\n")
				case eTypeUint16Array:
					cw.WriteLine("bb.PushUInt16((ushort)",vv.Name,".Count);")
					cw.WriteLine("for(int n=0; n < ",vv.Name,".Count; n++){")
					cw.DeepIn()
					cw.WriteLine("bb.PushUInt16(",vv.Name,"[n]);")
					cw.DeepOut()
					cw.WriteLine("}\n")
				case eTypeEnumArray:
					cw.WriteLine("bb.PushUInt16((ushort)",vv.Name,".Count);")
					cw.WriteLine("for(int n=0; n < ",vv.Name,".Count; n++){")
					cw.DeepIn()
					cw.WriteLine("bb.PushInt32((int)",vv.Name,"[n]);")
					cw.DeepOut()
					cw.WriteLine("}\n")
				case eTypeInt32Array:
					cw.WriteLine("bb.PushUInt16((ushort)",vv.Name,".Count);")
					cw.WriteLine("for(int n=0; n < ",vv.Name,".Count; n++){")
					cw.DeepIn()
					cw.WriteLine("bb.PushInt32(",vv.Name,"[n]);")
					cw.DeepOut()
					cw.WriteLine("}\n")
				case eTypeUint32Array:
					cw.WriteLine("bb.PushUInt16((ushort)",vv.Name,".Count);")
					cw.WriteLine("for(int n=0; n < ",vv.Name,".Count; n++){")
					cw.DeepIn()
					cw.WriteLine("bb.PushUInt32(",vv.Name,"[n]);")
					cw.DeepOut()
					cw.WriteLine("}\n")
				case eTypeInt64Array:
					cw.WriteLine("bb.PushUInt16((ushort)",vv.Name,".Count);")
					cw.WriteLine("for(int n=0; n < ",vv.Name,".Count; n++){")
					cw.DeepIn()
					cw.WriteLine("bb.PushInt64(",vv.Name,"[n]);")
					cw.DeepOut()
					cw.WriteLine("}\n")
				case eTypeUint64Array:
					cw.WriteLine("bb.PushUInt16((ushort)",vv.Name,".Count);")
					cw.WriteLine("for(int n=0; n < ",vv.Name,".Count; n++){")
					cw.DeepIn()
					cw.WriteLine("bb.PushUInt64(",vv.Name,"[n]);")
					cw.DeepOut()
					cw.WriteLine("}\n")
				case eTypeFloat32Array:
					cw.WriteLine("bb.PushUInt16((ushort)",vv.Name,".Count);")
					cw.WriteLine("for(int n=0; n < ",vv.Name,".Count; n++){")
					cw.DeepIn()
					cw.WriteLine("bb.PushFloat32(",vv.Name,"[n]);")
					cw.DeepOut()
					cw.WriteLine("}\n")
				case eTypeFloat64Array:
					cw.WriteLine("bb.PushUInt16((ushort)",vv.Name,".Count);")
					cw.WriteLine("for(int n=0; n < ",vv.Name,".Count; n++){")
					cw.DeepIn()
					cw.WriteLine("bb.PushFloat64(",vv.Name,"[n]);")
					cw.DeepOut()
					cw.WriteLine("}\n")
				case eTypeStringArray:
					cw.WriteLine("bb.PushUInt16((ushort)",vv.Name,".Count);")
					cw.WriteLine("for(int n=0; n < ",vv.Name,".Count; n++){")
					cw.DeepIn()
					cw.WriteLine("bb.PushString(",vv.Name,"[n]);")
					cw.DeepOut()
					cw.WriteLine("}\n")

				}
			}

			cw.WriteLine("return bb.Size();")
			cw.DeepOut()
			cw.WriteLine("}\n")

			//反序列化
			cw.WriteLine("public void deserialize(byte[] buff){")
			cw.DeepIn()
			cw.WriteLine("BuffParser bp = new BuffParser(buff);")
			cw.WriteLine("ushort _member_count = bp.PopUInt16();")
			cw.WriteLine("ushort _member_code = 0;")
			cw.WriteLine("byte _member_type = 0;")
			cw.WriteLine("for(var i=0;i<_member_count;i++)")
			cw.WriteLine("{")
			cw.DeepIn()

			cw.WriteLine("_member_code = bp.PopUInt16();")
			cw.WriteLine("_member_type = bp.PopUInt8();")
			cw.WriteLine("switch(_member_code){")
			cw.DeepIn()
			for _,vv := range v.Values {

				cw.WriteLine("case ",vv.Index,":")
				switch vv.Type{
				case eTypeInt8:
					cw.WriteLine(vv.Name+" = bp.PopInt8();break;")
				case eTypeUint8:
					cw.WriteLine(vv.Name+"  = bp.PopUInt8();break;")
				case eTypeInt16:
					cw.WriteLine(vv.Name+" = bp.PopInt16();break;")
				case eTypeUint16:
					cw.WriteLine(vv.Name+" = bp.PopUInt16();break;")
				case eTypeInt32:
					cw.WriteLine(vv.Name+" = bp.PopInt32();break;")
				case eTypeUint32:
					cw.WriteLine(vv.Name+" = bp.PopUInt32();break;")
				case eTypeFloat32:
					cw.WriteLine(vv.Name+" = bp.Float32();break;")
				case eTypeFloat64:
					cw.WriteLine(vv.Name+" = bp.Float64();break;")
				case eTypeString:
					cw.WriteLine(vv.Name+" = bp.PopString();break;")
				case eTypeEnum:
					cw.WriteLine(vv.Name+" = (",vv.BaseTypeName,")bp.PopInt32();break;")//TODO 检查值合法
				case eTypeInt8Array:
					cw.WriteLine("{")
					cw.DeepIn()
					cw.WriteLine("ushort arraylen = bp.PopUInt16();")
					cw.WriteLine("for (;arraylen>0;arraylen--){")
					cw.DeepIn()
					cw.WriteLine(vv.BaseTypeName+" element = bp.PopInt8();")
					cw.WriteLine(vv.Name+".Add(element);")
					cw.DeepOut()
					cw.WriteLine("}")
					cw.DeepOut()
					cw.WriteLine("}")
					cw.WriteLine("break;")
				case eTypeUint8Array:
					cw.WriteLine("{")
					cw.DeepIn()
					cw.WriteLine("ushort arraylen = bp.PopUInt16();")
					cw.WriteLine("for (;arraylen>0;arraylen--){")
					cw.DeepIn()
					cw.WriteLine(vv.BaseTypeName+" element = bp.PopUInt8();")
					cw.WriteLine(vv.Name+".Add(element);")
					cw.DeepOut()
					cw.WriteLine("}")
					cw.DeepOut()
					cw.WriteLine("}")
					cw.WriteLine("break;")
				case eTypeInt16Array:
					cw.WriteLine("{")
					cw.DeepIn()
					cw.WriteLine("ushort arraylen = bp.PopUInt16();")
					cw.WriteLine("for (;arraylen>0;arraylen--){")
					cw.DeepIn()
					cw.WriteLine(vv.BaseTypeName+" element = bp.PopInt16();")
					cw.WriteLine(vv.Name+".Add(element);")
					cw.DeepOut()
					cw.WriteLine("}")
					cw.DeepOut()
					cw.WriteLine("}")
					cw.WriteLine("break;")
				case eTypeUint16Array:
					cw.WriteLine("{")
					cw.DeepIn()
					cw.WriteLine("ushort arraylen = bp.PopUInt16();")
					cw.WriteLine("for (;arraylen>0;arraylen--){")
					cw.DeepIn()
					cw.WriteLine(vv.BaseTypeName+" element = bp.PopUInt16();")
					cw.WriteLine(vv.Name+".Add(element);")
					cw.DeepOut()
					cw.WriteLine("}")
					cw.DeepOut()
					cw.WriteLine("}")
					cw.WriteLine("break;")
				case eTypeInt32Array:
					cw.WriteLine("{")
					cw.DeepIn()
					cw.WriteLine("ushort arraylen = bp.PopUInt16();")
					cw.WriteLine("for (;arraylen>0;arraylen--){")
					cw.DeepIn()
					cw.WriteLine(vv.BaseTypeName+" element = bp.PopInt32();")
					cw.WriteLine(vv.Name+".Add(element);")
					cw.DeepOut()
					cw.WriteLine("}")
					cw.DeepOut()
					cw.WriteLine("}")
					cw.WriteLine("break;")
				case eTypeUint32Array:
					cw.WriteLine("{")
					cw.DeepIn()
					cw.WriteLine("ushort arraylen = bp.PopUInt16();")
					cw.WriteLine("for (;arraylen>0;arraylen--){")
					cw.DeepIn()
					cw.WriteLine(vv.BaseTypeName+" element = bp.PopUInt32();")
					cw.WriteLine(vv.Name+".Add(element);")
					cw.DeepOut()
					cw.WriteLine("}")
					cw.DeepOut()
					cw.WriteLine("}")
					cw.WriteLine("break;")
				case eTypeFloat32Array:
					cw.WriteLine("{")
					cw.DeepIn()
					cw.WriteLine("ushort arraylen = bp.PopUInt16();")
					cw.WriteLine("for (;arraylen>0;arraylen--){")
					cw.DeepIn()
					cw.WriteLine(vv.BaseTypeName+" element = bp.PopFloat32();")
					cw.WriteLine(vv.Name+".Add(element);")
					cw.DeepOut()
					cw.WriteLine("}")
					cw.DeepOut()
					cw.WriteLine("}")
					cw.WriteLine("break;")
				case eTypeFloat64Array:
					cw.WriteLine("{")
					cw.DeepIn()
					cw.WriteLine("ushort arraylen = bp.PopUInt16();")
					cw.WriteLine("for (;arraylen>0;arraylen--){")
					cw.DeepIn()
					cw.WriteLine(vv.BaseTypeName+" element = bp.PopFloat64();")
					cw.WriteLine(vv.Name+".Add(element);")
					cw.DeepOut()
					cw.WriteLine("}")
					cw.DeepOut()
					cw.WriteLine("}")
					cw.WriteLine("break;")
				case eTypeEnumArray:
					cw.WriteLine("{")
					cw.DeepIn()
					cw.WriteLine("ushort arraylen = bp.PopUInt16();")
					cw.WriteLine("for (;arraylen>0;arraylen--){")
					cw.DeepIn()
					cw.WriteLine(vv.BaseTypeName+" element = (",vv.BaseTypeName,")bp.PopInt32();")
					cw.WriteLine(vv.Name+".Add(element);")
					cw.DeepOut()
					cw.WriteLine("}")
					cw.DeepOut()
					cw.WriteLine("}")
					cw.WriteLine("break;")
				case eTypeStruct:
					cw.WriteLine("{")
					cw.DeepIn()
					cw.WriteLine("ushort elemlen = bp.PopUInt16();")
					cw.WriteLine("byte[] _buf = bp.PopBytes(elemlen);")
					cw.WriteLine(vv.BaseTypeName+" element = new ",vv.BaseTypeName,"();")
					cw.WriteLine("element.deserialize(_buf);")
					cw.DeepOut()
					cw.WriteLine("}")
					cw.WriteLine("break;")
				case eTypeStructArray:
					cw.WriteLine("{")
					cw.DeepIn()
					cw.WriteLine("ushort arraylen = bp.PopUInt16();")
					cw.WriteLine("for (;arraylen>0;arraylen--){")
					cw.DeepIn()
					cw.WriteLine("ushort elemlen = bp.PopUInt16();")
					cw.WriteLine("byte[] _buf = bp.PopBytes(elemlen);")
					cw.WriteLine(vv.BaseTypeName+" element = new ",vv.BaseTypeName,"();")
					cw.WriteLine("element.deserialize(_buf);")
					cw.WriteLine(vv.Name+".Add(element);")
					cw.DeepOut()
					cw.WriteLine("}")
					cw.DeepOut()
					cw.WriteLine("}")
					cw.WriteLine("break;")
				default:
					panic(errors.New("unknown type"))
				}

			}
			cw.WriteLine("default:")
			cw.WriteLine("bp.SkipType((eMemberType)_member_type);break;")
			cw.DeepOut()
			cw.WriteLine("}")//switch
			cw.DeepOut()
			cw.WriteLine("}")//for
			cw.DeepOut()

			cw.WriteLine("}\n")//deserialize
			cw.DeepOut()
			cw.WriteLine("}\n")//class
		}
	}

	//ns
	cw.DeepOut()
	cw.WriteLine("}")//namespace
}


func NewLexer() Lexer{
	l := &lexer{}
	l.Init()
	return l
}
