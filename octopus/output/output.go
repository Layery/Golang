package output

import (
	"fmt"
	"os"
	"sync"
)

// FileWriter 封装了文件写入和锁
type FileWriter struct {
	mu   sync.Mutex
	file *os.File
}

// NewFileWriter 创建新的 FileWriter
func NewFileWriter(filename string) (*FileWriter, error) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return nil, fmt.Errorf("打开文件失败: %w", err)
	}
	return &FileWriter{file: file}, nil
}

// Write 写入一行数据（线程安全）
func (fw *FileWriter) Write(title, url string) error {
	fw.mu.Lock()
	defer fw.mu.Unlock()
	line := fmt.Sprintf("<a style='margin-left:30px;margin-bottom:5px;' href='%s' target='_blank'>%s</a><br>\n", url, title)
	_, err := fw.file.Write([]byte(line))
	return err
}

// Close 关闭文件
func (fw *FileWriter) Close() error {
	fw.mu.Lock()
	defer fw.mu.Unlock()
	return fw.file.Close()
}
