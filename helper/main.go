package main

import (
	"fmt"
	"slices"
	"time"
)

type cost struct {
	day   int
	value float64
}

// 自己写的这个版本好low啊😅
func getResult(costList []cost) []float64 {
	result := []float64{}     // 存放最终结果
	temp := map[int]float64{} // 存放按照day聚合value后的数据
	keyList := []int{}        // 存放day值的切片, 排序后求出最大day
	for _, v := range costList {
		temp[v.day] += v.value
		if !slices.Contains(keyList, v.day) {
			keyList = append(keyList, v.day)
		}
	}

	// 按照day的大小排序, 取出最大的day, 并以这个day的值构造返回值
	slices.Sort(keyList)

	for i := 0; i <= keyList[len(keyList)-1]; i++ {
		if val, ok := temp[i]; ok {
			result = append(result, val)
		} else {
			result = append(result, 0.0)
		}
	}
	return result
}

// 这个方法是AI写的, 还是AI写的更精简啊😂
func aiGetCostByDay(input []cost) []float64 {
	if len(input) == 0 {
		return []float64{}
	}
	// 1. 找出最大的 day，用来决定输出切片的长度
	maxDay := 0
	for _, c := range input {
		if c.day > maxDay {
			maxDay = c.day
		}
	}

	// 2. 初始化结果切片，长度为 maxDay + 1 (因为 day 从 0 开始)
	// Go 会自动将切片里的所有元素初始化为 0.0
	result := make([]float64, maxDay+1)

	// 3. 遍历输入数据，按 day 进行聚合累加
	for _, c := range input {
		result[c.day] += c.value
	}

	return result
}

// 教程视频里这个方法更精简, 他没有先计算最大day的值, 而是当costByDay的长度不够了, 就append一个0来扩容
// 本质上其实是因为day的值和costByDay的索引下标一一对应, 所以内层for循环可以这么写,
// 但是当数据量很大的时候, 内层for循环, 会频繁的给切片扩容, 而扩容操作在底层是拷贝原切片到一个新的切片, 这样一来性能不高
func getCostByDay(costs []cost) []float64 {
	costByDay := []float64{}
	for i := 0; i < len(costs); i++ {
		cost := costs[i]
		for cost.day >= len(costByDay) {
			costByDay = append(costByDay, 0.0)
		}
		costByDay[cost.day] += cost.value
	}
	return costByDay
}

func LogRunTime(fn func(), funcName string) {
	start := time.Now()
	fn()
	fmt.Printf("函数 [%s] 耗时: %d毫秒\n", funcName, time.Since(start).Milliseconds())
}

func main() {

	// 写一个函数, 输入如下数据
	input := []cost{
		{9000000, 10},
		{0, 4.0},
		{1, 2.1},
		{1, 3.1},
		{5, 2.5},
		{8, 3.3},
		{8, 3.3},
	}

	// 输出如下结果, 这个结果等于是把cost中每个day的value求了下合, 然后按照day聚合, 中间缺失的day则展示0.0
	/**
	  []float64{
	      4.0,
	      5.2,
	      0.0,
	      0.0,
	      0.0,
	      2.5,
	  }
	*/

	LogRunTime(func() {
		getResult(input)
	}, "getResult")

	LogRunTime(func() {
		aiGetCostByDay(input)
	}, "aiGetCostByDay")

	LogRunTime(func() {
		getCostByDay(input)
	}, "getCostByDay")
}
