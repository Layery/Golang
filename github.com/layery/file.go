package main

import (
	"fmt"
	"net/http"
)

func main()  {
	http.Handle("/", http.FileServer(http.Dir("E:/")))
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		fmt.Println(err)
	}
}