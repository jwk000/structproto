# structproto
类似protobuf的二进制兼容的结构体序列化工具，可生成C++、C#代码。

##支持注释 //xxx
##类似protobuf语法
需要写一个xxx.struct文件描述结构体
支持数据类型 int8 int16 int32 int64 uint8 uint16 uint32 uint64 float32 float64 string enum
支持数组 类型后跟一个+号表示repeated 字段：比如 itemIds int32+ = 1; 表示一个int list.
支持结构体嵌套
二进制兼容的约定：和pb一样，数字id不能复用，可以删除，新增

##结构体字段的处理规则

1.最外部结构体处理（就是用户描述的那些结构体）
uint16成员个数n
(uint16编号 uint8类型 数据) n组

2字符串的处理
不带结尾0，按utf8编码
uint16长度 字符串内容

3.内部结构体处理(就是作为数据成员的结构体)
uint16长度 结构体数据

4.数组的处理
uint16元素个数 元素数据

##TODO
//  支持属性 "属性1=值1 属性2=值2"
//  范围检查"min=1 max=100"
//  默认值"default=0"
//  对string字段有长度检查"maxstrlen=32"
//  对repeated字段有数量最少最大检查"mincount=1 maxcount=10"
//  对枚举字段有值范围检查
//使用message标记可以用作消息结构体发送，生成协议号，自动调用callback
//推荐的命名方式：前缀_协议名称
//使用者直接发送结构体，不需要定义协议号，代码生成协议号
//协议头、加密在上层做，和协议内容无关
//生成代码包含函数：MinSize() 返回结构体最小长度用来验证消息大小 MaxSize()同理



##下面是测试用的c2s.struct文件，请放在bin/目录下
```
enum ePlayerState{
    offline=0
    online=1
    inteam=2
    ingame=3
}

struct s2c_player_info{
    playerid int = 1
    playername string = 2
    playerstate ePlayerState = 3
}

struct s2c_all_player_info{
    infos s2c_player_info+ = 1
}

struct msg_stone{
    keyid   int = 1
    stoneid int = 2
}

struct s2c_refresh_stone{
    stones msg_stone+ = 1
}
```
