package np

import (
	"errors"
	"fmt"
	//	"bytes"
	"crypto/md5"
	"encoding/binary"
	"io"
	"time"

	"net"
	//	"bufio"
	//"log"
)

const (
	// 用于指定帧的大小，不包括帧头
	MaxBodySize = 1024
	// 当发送队列为空时，指定每次检查队列状况的时间间隔
	WriteInterval = 80
	// 当读到EOF时，指定每次重读的时间间隔
	ReadInterval = 80
	//
	MaxInvalidPackageNum = 2000
)

var (
	ChanClosedErr   = errors.New("netchan has been closed!")
	SocketClosedErr = errors.New("socket has been closed!")
	ChanNotExistErr = errors.New("netchan does not exist!")
	DataConvertErr  = errors.New("data convert failed!")
)

type NetChanSet map[[md5.Size]byte]*NetChan
type Socket struct {
	conn              *net.TCPConn
	netchans          NetChanSet
	isClose           bool
	invalidPackageNum uint32
}

// 新建一个管道管理器
func New(conn *net.TCPConn) *Socket {
	s := new(Socket)
	s.conn = conn
	s.isClose = false
	s.netchans = make(NetChanSet)

	s.conn.SetNoDelay(false)
	// 用于通知管道关闭的管道
	c := NewBlockChan(md5.Size, new(ChanidConverter))
	callback := func(i interface{}) bool {
		// 关闭对应管道
		if v, ok := s.netchans[i.([md5.Size]byte)]; ok {
			v.Close(ChanClosedErr)
			// TODO delete(s.netchans, i.([md5.Size]byte))
		}
		return true
	}
	c.SetReceiveCallBack(callback)
	s.RegisterChannel("CloseChan", c)
	// TODO 心跳包管道

	go s.readDaemon()  //启动轮寻读协程，接收数据
	go s.writeDaemon() //启动轮寻写协程，发送数据
	return s
}

// 关闭链接
func (s *Socket) Close() error {
	s.isClose = true
	for _, netchan := range s.netchans {
		netchan.Close(SocketClosedErr)
	}
	err := s.conn.Close()
	return err
}

// 注册一个管道
func (s *Socket) RegisterChannel(name string, netchan *NetChan) {
	// TODO
	s.netchans[md5.Sum([]byte(name))] = netchan
}

// 注销一个管道
func (s *Socket) UnregisterChannel(name string) {
	if s.isClose {
		return
	}
	s.netchans[md5.Sum([]byte(name))].Close(ChanClosedErr)
	delete(s.netchans, md5.Sum([]byte(name)))
}

// 用于轮讯另一端发过来的数据
func (s *Socket) readDaemon() {

	temp := make([]byte, headLength)

	for {
		if s.isClose {
			return
		}

		err := s.read(&temp)
		if err != nil {
			fmt.Println("readDaemon() read head failed: " + err.Error())
			s.Close()
			return
		}

		h := createHeadFrom(&temp)
		if h.bodySize > MaxBodySize || h.bodySize == 0 {
			fmt.Println("readDaemon() error: body too large")
			s.Close()
			return
		}

		body := make([]byte, h.bodySize)
		err = s.read(&body)
		if err != nil {
			fmt.Println("readDaemon() read body failed: " + err.Error())
			s.Close()
			return
		}

		// 转发到不同的管道
		netchan, ok := s.netchans[h.chanId]
		if ok == true {
			if netchan.receive(body) != nil {
				// 通知另一端管道被关闭
				s.Put("CloseChan", h.chanId)
				// 关闭管道
				netchan.Close(err)
				// delete(s.netchans, h.chanId)
			}
		} else {
			// 记录下无对应管道的包的数量，当大于规定数量时关闭socket
			s.invalidPackageNum++
			if s.invalidPackageNum > MaxInvalidPackageNum {
				s.Close()
				return
			}
			continue
		}
	}

}

func (s *Socket) read(b *[]byte) error {
	for readSize := 0; readSize < len(*b); {
		count, err := s.conn.Read((*b)[readSize:])
		if err != nil && err != io.EOF {
			return err
		}
		if err == io.EOF {
			time.Sleep(ReadInterval * time.Millisecond)
		}
		readSize += count
	}
	return nil
}

// 用于轮讯发送数据
func (s *Socket) writeDaemon() {

	var noWrite bool

	for {
		if s.isClose {
			return
		}

		noWrite = true
		for chanid, netchan := range s.netchans {

			body := netchan.getOneSendFrame()

			if body == nil {
				continue
			}
			// 构造frame
			h := new(head)
			h.chanId = chanid
			h.bodySize = uint64(len(*body))
			buf := make([]byte, headLength+h.bodySize)
			copy(buf, h.toByteArray())
			copy(buf[headLength:], *body)
			// 发送数据
			err := s.write(&buf)
			if err != nil {
				fmt.Println("writeDaemon() write failed: " + err.Error())
				s.Close()
				return
			}
			noWrite = false
		}
		if noWrite {
			time.Sleep(WriteInterval * time.Millisecond)
		}

	}
}

func (s *Socket) write(b *[]byte) error {
	for writeSize := 0; writeSize < len(*b); {
		count, err := s.conn.Write((*b)[writeSize:])
		if err != nil {
			return err
		}
		writeSize += count
	}
	return nil
}

// 从一个管道读取数据, 该函数是线程安全的，支持并发访问
func (s *Socket) Get(chaname string) (interface{}, error) {
	if s.isClose {
		return nil, SocketClosedErr
	}

	v, ok := s.netchans[md5.Sum([]byte(chaname))]
	if ok {
		return v.get()
	} else {
		return nil, ChanNotExistErr
	}
}

// 将数据写入一个管道，该函数是线程安全的，支持并发访问
func (s *Socket) Put(chaname string, i interface{}) error {
	if s.isClose {
		return SocketClosedErr
	}

	v, ok := s.netchans[md5.Sum([]byte(chaname))]
	if ok {
		return v.put(&i)
	} else {
		return ChanNotExistErr
	}
}

const (
	// 定义帧头的长度
	headLength = md5.Size + 8 // 8是uint64的长度
)

// 内部数据帧的帧头
type head struct {
	chanId   [md5.Size]byte
	bodySize uint64
}

func createHeadFrom(b *[]byte) *head {
	h := new(head)
	// 拷贝chanid字段
	copy(h.chanId[:], (*b)[:md5.Size])
	// 拷贝bodysize字段
	h.bodySize = binary.LittleEndian.Uint64((*b)[md5.Size:])

	return h
}

func (h *head) toByteArray() []byte {
	temp := make([]byte, headLength)
	// 拷贝chanid字段
	copy(temp, h.chanId[:])
	// 拷贝bodysize字段
	binary.LittleEndian.PutUint64(temp[md5.Size:], h.bodySize)

	return temp
}
