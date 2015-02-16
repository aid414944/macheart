[TOC]
#NP
np取自new protocol之意，使用它可以在一个TCP链接上创建多个数据传输管道；管道支持并发的发送和接收操作。
##使用
下面是一个简单的使用演示：
```go
// 初始化
socket := np.New(conn) // conn 是一个net.TCPConn类型的指针
netchan := np.NewStreamChan([]byte{0, 0}, new(np.StringConverter))// 新建一个流式管道
socket.RegisterChannel("text", netchan)
// 发送数据
socket.Put("text", "hello world!")
// 请求数据
str, err := socket.Get("text")
// 接收数据
callback := func(object interface{}) bool {
	// handle data
	return true
}
netchan.SetReceiveCallBack(callback)
```
np支持两种类型的管道：
* 流式：调用np.NewStreamChan创建
* 块式：调用np.NewBlockChan创建

对于字符串这种长度不定的数据类型，需要使用流式管道。
##扩展
一个管道只能传输一种类型的数据，具体是那种，由管道的构造函数的输入参数决定；对于自定义的数据类型，只需要实现相应的np.Converter接口，并传递给管道的构造函数即可。
```go
type Converter interface {
	ToObject(data *[]byte) (interface{}, error)
	ToArrayByte(object *interface{}) ([]byte, error)
}
```