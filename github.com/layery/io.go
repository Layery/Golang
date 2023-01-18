package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
)

func DD(data ...interface{}) { // 声明可变参数函数时，需要在参数列表的最后一个参数类型之前加上省略符号“...”，这表示该函数会接收任意数量的该类型参数。
	d, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}

	var out bytes.Buffer
	err = json.Indent(&out, d, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(out.String())
}

// 从标准输入读入数据, 并直接写入到文件中
func WriteFileFromBufio() {
	fd, err := os.OpenFile("./io.log", os.O_CREATE, 0777)
	if err != nil {
		log.Fatal(err)
	}

	defer fd.Close()

	reader := bufio.NewReader(os.Stdin)

	p := make([]byte, 3)
	for {
		n, err := reader.Read(p)
		if err != nil {
			if err == io.EOF {
				log.Fatal("read over")
				break
			}
			log.Fatal(err)
		}

		// 将读到的内容, 实时写入文件中去
		if n > 0 {
			_, writeErr := fd.WriteString(string(p[0:n]))
			if writeErr != nil {
				log.Fatal(writeErr)
			}

		}
		fmt.Println(string(p[0:n]))
	}
}

func readDir(path string) {
	rs, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	for key, value := range rs {
		_ = key
		if value.IsDir() {
			readDir(path + string(os.PathSeparator) + value.Name())
		}

		log.Println(value.Name())
	}

}

func main() {

	// go WriteFileFromBufio()

	go readDir(".")

	select {}

}
