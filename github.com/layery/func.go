package main

import (
	"fmt"
)

func add(a int, b int) (int, int) {
	return a + b, a - b
}

func main() {
	// go的匿名函数

	sum, _ := add(3, 4)

	fmt.Println(sum)

	for i := 0; i < 5; i++ {
		defer fmt.Printf("%d ", i)
	}

}
