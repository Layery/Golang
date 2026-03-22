package main

import (
	"fmt"
	"github.com/gookit/goutil/dump"
	"sync"
)

type Goods struct {
	mu    sync.Mutex
	goods map[string]interface{}
}

func NewGoods() *Goods {
	return &Goods{
		goods: map[string]interface{}{
			"name":  "矿泉水",
			"stock": 100,
			"sales": 0,
		},
	}
}

// 尝试购买商品
func (g *Goods) TryPurchase() bool {
	g.mu.Lock()
	defer g.mu.Unlock()

	stock := g.goods["stock"].(int)
	if stock > 0 { // 检查库存是否充足
		g.goods["stock"] = stock - 1
		g.goods["sales"] = g.goods["sales"].(int) + 1
		return true
	}
	return false
}

// 获取当前商品信息
func (g *Goods) GetGoodsInfo() map[string]interface{} {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.copyGoods()
}

// 深拷贝商品信息，防止外部直接修改
func (g *Goods) copyGoods() map[string]interface{} {
	copied := make(map[string]interface{})
	for k, v := range g.goods {
		copied[k] = v
	}
	return copied
}

func map_bingfa_unsafe() {
	goods := NewGoods()
	var wg sync.WaitGroup

	// 模拟200个用户并发抢购
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			if goods.TryPurchase() {
				fmt.Printf("User %d 抢购成功了\n", id)
			} else {
				fmt.Printf("User %d 抢购失败\n", id)
			}
		}(i)
	}

	// 等待所有用户完成抢购
	wg.Wait()

	// 打印最终的商品信息
	info := goods.GetGoodsInfo()
	fmt.Printf("Final Goods Info: %+v\n", info)

}
func main() {
	// 模拟map 并发不安全
	//map_bingfa_unsafe()

	/**
	关于golang中map的这种古怪的特性有这样几个观点：
		1）map作为一个封装好的数据结构，由于它底层可能会由于数据扩张而进行迁移，所以拒绝直接寻址，避免产生野指针；
		2）map中的key在不存在的时候，赋值语句其实会进行新的k-v值的插入，所以拒绝直接寻址结构体内的字段，以防结构体不存在的时候可能造成的错误；
		3）这可能和map的并发不安全性相关

		x = y 这种赋值的方式，你必须知道 x的地址，然后才能把值 y 赋给 x。
		但 go 中的 map 的 value 本身是不可寻址的，因为 map 的扩容的时候，可能要做 key/val pair迁移
		value 本身地址是会改变的
		不支持寻址的话又怎么能赋值呢
	*/
	map_value_is_struct()
}

func map_value_is_struct() {
	type Student struct {
		Name string
	}
	// 声明一个map, 值为Student结构体
	list := map[string]Student{}
	// 给map的stu赋值为一个Student结构体
	list["stu"] = Student{Name: "llf"}

	// 此时是无法修改stu的值的Name属性的, 编译不过去
	//list["stu"].Name = "test"

	// 那么如何修改呢? 考虑如下两种方式:

	// 方式1: 将stu的值重新赋值一个临时变量, 虽然能实现修改结构体中的字段, 但是引发了一次值拷贝, 性能不好
	temp := list["stu"]
	temp.Name = "test"
	list["stu"] = temp
	fmt.Println(list)

	// 方式2: 直接将stu的结构体, 定义为引用类型, 则可直接修改
	list2 := map[string]*Student{}
	list2["stu"] = &Student{Name: "llf"}
	dump.Println(list2)
	list2["stu"].Name = "test"

	dump.P(list2)
}
