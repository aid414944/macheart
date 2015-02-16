package np

import (
	"fmt"
	"sync"
	//"time"
	"container/list"
)

type DataHandler interface {
	Load([]byte) bool // 当加载了足够多的数据能够进行转换时，函数返回true
	ToObject() (interface{}, error)

	Cut(*interface{}) error      //第二个参数是每片大小
	GetOneSlice() ([]byte, bool) // 如果没有数据，第二个返回值为假
	HasSlice() bool
}

type Converter interface {
	ToObject(data *[]byte) (interface{}, error)
	ToArrayByte(object *interface{}) ([]byte, error)
}

type NetChan struct {
	readTaskList    *list.List // 读取队列
	readTaskLock    sync.Mutex
	writeTaskList   *list.List // 写入队列
	writeTaskLock   sync.Mutex
	receiveCallBack func(interface{}) bool
	dataHandler     DataHandler
	isClosed        bool
}

// 新建一个流式管道。endFlag用于指定流的结束标记，请确保流中的数据不包含endFlag
func NewStreamChan(endFlag []byte, converter Converter) *NetChan {
	nc := NewChan(NewStreamHandler(endFlag, converter))
	return nc
}

// 新建一个块式管道。size用于指定块的大小，请确保它与converter转换出来的[]byte长度一致
func NewBlockChan(size uint, converter Converter) *NetChan {
	nc := NewChan(NewBlockHandler(size, converter))
	return nc
}

// 新建一个管道。
func NewChan(dataHandler DataHandler) *NetChan {
	nc := new(NetChan)
	nc.readTaskList = list.New()
	nc.writeTaskList = list.New()
	nc.dataHandler = dataHandler
	nc.isClosed = false
	return nc
}

// 关闭管道，所有等待中的任务将会退出，阻塞的函数将恢复执行。
// 参数err用于给调用者报告错误
func (nc *NetChan) Close(err error) {

	// 上锁
	nc.readTaskLock.Lock()
	nc.writeTaskLock.Lock()
	defer nc.readTaskLock.Unlock()
	defer nc.writeTaskLock.Unlock()

	if nc.isClosed == true {
		return
	}
	nc.isClosed = true

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("closeChan() error:", err)
		}
	}()

	// 释放接收队列里的任务
	e := nc.readTaskList.Front()
	for e != nil {
		// 恢复执行
		t := e.Value.(*task)
		*((*t).Data) = nil
		(*t).Error <- err
		// 删除任务
		temp := e.Next()
		nc.readTaskList.Remove(e)
		e = temp
	}
	// 释放发送队列里的任务
	e = nc.writeTaskList.Front()
	for e != nil {
		// 恢复执行
		t := e.Value.(*task)
		*((*t).Data) = nil
		(*t).Error <- err
		// 删除任务
		temp := e.Next()
		nc.writeTaskList.Remove(e)
		e = temp
	}

}

// 设置接收数据的回调函数，若回调返回true，则数据被消耗。
// 注意！不要在回调函数中执行太长时间的操作，否则会阻塞后续数据的接收和传递。
func (nc *NetChan) SetReceiveCallBack(callback func(interface{}) bool) {
	nc.receiveCallBack = callback
}

// 接收来自socket的数据，数据将会在此通过DataHandler构造为对象并分发出去。若管道的读取队列为空，数据会被抛弃。
// data的长度不会大于MaxBodySize。
// 返回nil表示数据被正确处理，没有产生错误。
// 返回非nil表示可能在数据转换过程中发生错误，这种状况通常是由于接收到不合法的数据造成的。
func (nc *NetChan) receive(data []byte) error {

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("NetChan receice() error:", err)
		}
	}()

	if ok := nc.dataHandler.Load(data); !ok {
		return nil
	}
	object, err := nc.dataHandler.ToObject()
	if err != nil {
		// TODO
		// 转换对象失败，说明接收到了不合法的数据
		// 应该在此关闭管道，释放队列中的任务
		// 设置关闭标志，不再接受读写任务
		// 下面的代码未测试
		// nc.Close(err) 这条代码已被移至np.go
		return err
	}

	if nc.receiveCallBack != nil {
		if nc.receiveCallBack(object) {
			return nil
		}
	}

	// 上锁
	nc.readTaskLock.Lock()
	defer nc.readTaskLock.Unlock()

	length := nc.readTaskList.Len()
	if length == 0 {
		return nil
	}

	e := nc.readTaskList.Front()
	for e != nil {
		t := e.Value.(*task)
		*((*t).Data) = object
		(*t).Error <- nil
		// 释放接收任务
		close((*t).Error)
		temp := e.Next()
		nc.readTaskList.Remove(e)
		e = temp
	}
	return nil
}

// 从管道中接收一个对象并返回，该函数是阻塞的。
func (nc *NetChan) get() (object interface{}, err error) {

	c := make(chan error)
	t := newTask(c, &object)
	nc.readTaskLock.Lock()
	if nc.isClosed {
		close(c)
		nc.readTaskLock.Unlock()
		return nil, ChanClosedErr
	}
	nc.readTaskList.PushBack(t)
	nc.readTaskLock.Unlock()

	err = <-c
	return
}

// 从发送队列的第一个任务里，抽取一帧待发送数据，并返回。
// 如果发送队列为空，则返回nil
func (nc *NetChan) getOneSendFrame() *[]byte {
	// 上锁
	nc.writeTaskLock.Lock()
	defer nc.writeTaskLock.Unlock()

	if nc.isClosed {
		return nil
	}

	length := nc.writeTaskList.Len()
	if length == 0 {
		return nil
	}

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("NetChan GetOneSendFrame() error:", err)
		}
	}()

	e := nc.writeTaskList.Front()
	t := e.Value.(*task)

	var buf []byte
	var hasNext bool
	if !nc.dataHandler.HasSlice() {
		err := nc.dataHandler.Cut((*t).Data)
		if err != nil {
			// 转换为[]byte失败
			// 释放任务，传递错误信息并返回
			(*t).Error <- err
			// 释放发送任务
			close((*t).Error)
			nc.writeTaskList.Remove(e)
			return nil
		}
	}

	buf, hasNext = nc.dataHandler.GetOneSlice()
	if !hasNext {
		(*t).Error <- nil
		// 释放发送任务
		close((*t).Error)
		nc.writeTaskList.Remove(e)
	}

	return &buf

}

// 通过管道发送一个对象
func (nc *NetChan) put(data *interface{}) (err error) {

	c := make(chan error)
	t := newTask(c, data)
	nc.writeTaskLock.Lock()
	if nc.isClosed {
		close(c)
		nc.writeTaskLock.Unlock()
		return ChanClosedErr
	}
	nc.writeTaskList.PushBack(t)
	nc.writeTaskLock.Unlock()
	err = <-c

	return
}

type task struct {
	Error chan error
	Data  *interface{}
}

func newTask(chanError chan error, data *interface{}) *task {
	t := new(task)
	t.Error = chanError
	t.Data = data
	return t
}
