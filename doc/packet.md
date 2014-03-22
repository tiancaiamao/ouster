# 网络封包

## 服务端使用接口

packet对外提供的接口为

	func Marshal(v interface{}) ([]byte, error)
	func Unmarshal(data []byte, v interface{}) error
	
Unmarshal返回一个interface，这个interface其实是一个packet，通过

	switch t.(type) {
		case *packet.LoginPacket
		case *packet.SelectCharactorPacket
	}

或者
	
	pkt.(*packet.LoginPacket)

类似方法来判断具体是哪个包。

## 实现

在网络上传输的包的结构是: 包大小|包的id | 包数据。Parse接收的数据块是除去包大小后的部分，首先是以一个包id开头的，然后是具体内容。

通过id可以区分出是哪个包，packet模块中有一张表，做id到reflect.Type的映射表。通过reflect.Type里的信息对包进行解释，构造出具体的数据结构。这样做的优势是包可以做到非常小，因为不必把类型信息放进包里。

包内容格式决定使用msgpack库。善良的国际友人给我们造了那么多轮子，如果我不使用，他们会很伤心的。

比如一个结构体：

	type Test struct {
		Name string
		Num uint32
	}

存储形式就是：Test包对应的id | msgpack打包的字符串值 | msgpack打包的 uint32值

## 客户端

客户端那边由于C语言是没有reflect的，采用代码生成的方式。类似protobuf的做法，proto就是Go语言的packet包中定义的那些结构体。

用Go语言写个小工具，遍历packet中的各种包的reflect.Type结构，生成对应的读写包的C代码。但是由于没有反射，C那边代码会写得丑一些。

接口： 

	int Unmarshal(char *data, struct Packet * pkt,);
	int Marshal(struct Packet* pkt, struct Slice *buf); 

struct Packet {
	uint16_t Id;
	char Data;
};

每种具体的struct Packet中都有一个uint16_t的Id来表明这是一个什么样的packet。有了这个id以后就可以做强制类型转换得到具体的Packet结构体。这些都应该是代码生成的。

## 代码生成

假设对于封包

	type Test struct {
		Name string
		Num uint32
	}

在Go那边是利用id+类型反射信息处理的，而C这边没有反射信息，使用代码生成的方式。

生成的代码应该是这样子的:

#define PTest 34

switch(pkt->Id) {
	case PTest:
		ReadString(buf, &((struct CharactorInfoPacket*)pkt)->Name);
		ReadS32(buf, &((struct CharactorInfoPacket*)pkt)->Level);
		break;
}