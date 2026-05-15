package main

import "fmt"

type TNumber any {
	int | float64 | string
}

func getSumByAny[Number TNumber](n1, n2 Number) Number {
	return n1 * n2
}
func main() {
	resp := getSumByAny(1.3, 3)
	fmt.Println(resp)
}
