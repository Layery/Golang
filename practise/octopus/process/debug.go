package process

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/urfave/cli/v2"
)

func RunDebug(ctx *cli.Context) error {
	file, err := os.OpenFile("./output/dist/debug.log", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("打开日志文件失败: %w", err)
	}
	defer file.Close()

	mw := io.MultiWriter(os.Stdout, file)

	demoWriter(mw)
	demoReader(mw)
	demoTeeReader(mw)
	demoCopy(mw)

	fmt.Fprintln(mw, "=== 全部演示结束 ===")
	return nil
}

// demoWriter 演示 io.Writer 的各种写入方式
func demoWriter(w io.Writer) {
	fmt.Fprintln(w, "==============================")
	fmt.Fprintln(w, "【io.Writer 演示】")
	fmt.Fprintln(w, "==============================")
	fmt.Fprintln(w)

	// 1. fmt.Fprintln — 写入字符串并追加换行
	fmt.Fprintln(w, "1. fmt.Fprintln: 写入字符串")

	// 2. fmt.Fprintf — 格式化写入
	fmt.Fprintf(w, "2. fmt.Fprintf:  PI = %.10f\n", 3.1415926535)

	// 3. io.Writer.Write() — 最底层的方法，写入 []byte
	data := []byte("3. Write([]byte): 这是原始字节数据\n")
	n, _ := w.Write(data)
	fmt.Fprintf(w, "   写入了 %d 字节\n", n)

	// 4. 使用 bytes.Buffer 作为 Writer，写入后再取出
	var buf bytes.Buffer
	fmt.Fprint(&buf, "4. bytes.Buffer: ")
	buf.WriteString("先写入 Buffer ")
	buf.WriteString("再写入更多内容\n")
	buf.WriteTo(w) // Buffer 实现了 io.WriterTo，将内容写入 w

	fmt.Fprintln(w)
}

// demoReader 演示 io.Reader 的各种读取方式
func demoReader(w io.Writer) {
	fmt.Fprintln(w, "==============================")
	fmt.Fprintln(w, "【io.Reader 演示】")
	fmt.Fprintln(w, "==============================")
	fmt.Fprintln(w)

	// 1. strings.NewReader — 从字符串创建 Reader
	reader := strings.NewReader("Hello, io.Reader!")
	buf := make([]byte, 1024)
	n, _ := reader.Read(buf)
	fmt.Fprintf(w, "1. strings.NewReader + Read: %s\n", buf[:n])

	// 2. 重复读取（Seek 回到开头）
	reader.Seek(0, io.SeekStart)
	n, _ = io.ReadFull(reader, buf)
	fmt.Fprintf(w, "2. ReadFull (读取全部):     %s\n", buf[:n])

	// 3. io.ReadAll — 一次性读取全部内容
	reader.Seek(0, io.SeekStart)
	all, _ := io.ReadAll(reader)
	fmt.Fprintf(w, "3. ReadAll:                 %s\n", all)

	// 4. bytes.NewReader — 从 []byte 创建 Reader
	byteData := []byte("从字节切片读取数据")
	byteReader := bytes.NewReader(byteData)
	chunk := make([]byte, 30)
	n, _ = byteReader.Read(chunk)
	fmt.Fprintf(w, "4. bytes.NewReader + Read:  %s (共%d字节)\n", chunk[:n], len(byteData))

	// 5. 分块读取 — 模拟大文件分块处理
	fmt.Fprintln(w, "5. 分块读取演示:")
	longReader := strings.NewReader("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	smallBuf := make([]byte, 8)
	for {
		n, err := longReader.Read(smallBuf)
		if n > 0 {
			fmt.Fprintf(w, "   读取 %d 字节: [%s]\n", n, smallBuf[:n])
		}
		if err == io.EOF {
			break
		}
	}

	fmt.Fprintln(w)
}

// demoTeeReader 演示 io.TeeReader — 读取的同时自动写入
func demoTeeReader(w io.Writer) {
	fmt.Fprintln(w, "==============================")
	fmt.Fprintln(w, "【io.TeeReader 演示】")
	fmt.Fprintln(w, "==============================")
	fmt.Fprintln(w)

	src := strings.NewReader("TeeReader: 读一份，写两份")
	var mirror bytes.Buffer // 镜像 Writer，自动接收读过的内容

	// TeeReader 返回一个 Reader，每次读取时会自动把内容写入 mirror
	tee := io.TeeReader(src, &mirror)

	buf := make([]byte, 1024)
	n, _ := tee.Read(buf)
	fmt.Fprintf(w, "读取到:   %s\n", buf[:n])
	fmt.Fprintf(w, "镜像中:   %s\n", mirror.String())
	fmt.Fprintln(w)
}

// demoCopy 演示 io.Copy — Reader → Writer 的管道
func demoCopy(w io.Writer) {
	fmt.Fprintln(w, "==============================")
	fmt.Fprintln(w, "【io.Copy 演示】")
	fmt.Fprintln(w, "==============================")
	fmt.Fprintln(w)

	// 1. 基本 Copy: strings.Reader → bytes.Buffer
	src := strings.NewReader("io.Copy: 数据从 Reader 流向 Writer")
	var dst bytes.Buffer
	io.Copy(&dst, src)
	fmt.Fprintf(w, "1. Copy 结果: %s\n", dst.String())

	// 2. CopyN: 只复制指定字节数
	src.Seek(0, io.SeekStart)
	dst.Reset()
	io.CopyN(&dst, src, 8)
	fmt.Fprintf(w, "2. CopyN (8字节): %s\n", dst.String())

	// 3. CopyN 到多 Writer
	src.Seek(0, io.SeekStart)
	dst.Reset()
	multiDst := io.MultiWriter(&dst, w) // 同时写入 buffer 和终端
	fmt.Fprint(w, "3. CopyN → MultiWriter: ")
	io.CopyN(multiDst, src, 18) // 只复制前18字节 "io.Copy: 数据从"
	fmt.Fprintf(w, "\n   buffer 中也有: %s\n", dst.String())

	fmt.Fprintln(w)
}
