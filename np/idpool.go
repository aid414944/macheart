package np

import (
	"container/list"
	"sort"
)

// 已知问题：1.未加锁，多线程操作下会破坏数据
//		   2.源码未测试

type IDPool struct {
	upLimit     uint64
	downLimit   uint64
	nextID      uint64
	recycleList *list.List
}

// 新建一个ID池，downLimit指定ID的下限，upLimit指定ID的上限，downLimit必须小于upLimit
func NewIDPool(downLimit, upLimit uint64) *IDPool {
	idpool := new(IDPool)
	idpool.upLimit = upLimit
	idpool.downLimit = downLimit
	idpool.nextID = downLimit
	idpool.recycleList = list.New()
	return idpool
}

// 分配id
func (idpool *IDPool) Get() (uint64, bool) {
	if idpool.recycleList.Len() == 0 {
		if idpool.nextID == idpool.upLimit+1 {
			return 0, false
		}
		temp := idpool.nextID
		idpool.nextID += 1
		return temp, true
	}
	e := idpool.recycleList.Front()
	return idpool.recycleList.Remove(e).(uint64), true
}

// 回收id
func (idpool *IDPool) Recycle(id uint64) {
	idpool.recycleList.PushBack(id)
}

// 整理id池
func (idpool *IDPool) Tidy() {
	// 将recycleList转换为uint64Slice并排序
	array := make(uint64Slice, 0)
	e := idpool.recycleList.Front()
	for e != nil {
		array = append(array, e.Value.(uint64))
		e = e.Next()
	}
	sort.Sort(array)
	//
	var j int
	for j = len(array) - 1; j >= 0; j-- {
		if array[j] == idpool.nextID-1 {
			idpool.nextID -= 1
			continue
		} else {
			break
		}
	}
	//
	idpool.recycleList.Init()
	for i := 0; i <= j; i++ {
		idpool.recycleList.PushBack(array[i])
	}
}

type uint64Slice []uint64

func (us uint64Slice) Len() int {
	return len(us)
}

// 如果i索引的数据小于j索引的数据，返回true，不会调用
// 下面的Swap()，即数据升序排序。
func (us uint64Slice) Less(i, j int) bool {
	return us[i] < us[j]
}

func (us uint64Slice) Swap(i, j int) {
	us[i], us[j] = us[j], us[i] // 看了golang源码才知道有这么牛逼的写法
}
