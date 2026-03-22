package main

import "fmt"

func main() {
	arr := []int{1, 2, 3, 5, 6, 7}
	for i, i2 := range arr {
		defer func() {
			ii := i
			vv := i2
			fmt.Println(ii, vv)
		}()
	}

}
