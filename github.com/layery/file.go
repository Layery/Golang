package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.Handle("/", http.FileServer(http.Dir("/home/layery/")))
	err := http.ListenAndServe(":8181", nil)
	if err != nil {
		fmt.Println(err)
	}
}
