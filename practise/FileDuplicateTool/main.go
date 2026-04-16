package main

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/gookit/goutil/dump"
)

var fileTypeMap = map[string]map[string]bool{
	"music": {
		".mp3":  true,
		".flac": true,
		".wav":  true,
		".m4a":  true,
		".aac":  true,
		".ogg":  true,
		".wma":  true,
		".ape":  true,
		".alac": true,
		".amr":  true,
	},
	"video": {
		".mp4": true,
		".avi": true,
		".mov": true,
	},
}

func isMatchFileType(file os.DirEntry, fileType string) bool {
	ext := filepath.Ext(file.Name())
	// fileType = "uuuu"
	return fileTypeMap[fileType][ext]
}

func getFileInfo(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("无法打开文件: %s, 错误: %v", filePath, err)
	}
	defer file.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return "", fmt.Errorf("读取文件失败: %s, 错误: %v", filePath, err)
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}

func listDir(path string, fileType string, depth int) ([]string, error) {
	if depth > 100 {
		return nil, fmt.Errorf("递归深度超过限制: %s", path)
	}

	_, err := os.Stat(path)
	if errors.Is(err, fs.ErrNotExist) {
		return nil, fmt.Errorf("目录不存在: %s", path)
	}
	if err != nil {
		return nil, fmt.Errorf("无法访问路径: %s, 错误: %v", path, err)
	}

	f, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("读取目录失败: %s, 错误: %v", path, err)
	}

	result := make([]string, 0)

	for _, file := range f {
		_, err := file.Info()
		if err != nil {
			fmt.Printf("无法获取文件信息: %s \n", file.Name())
			continue
		}

		fullPath := filepath.Join(path, file.Name())

		if file.IsDir() {
			subFiles, err := listDir(fullPath, "music", depth+1)
			if err != nil {
				return nil, fmt.Errorf("递归错误: %v", err)
			}
			result = append(result, subFiles...)
		} else {
			simple_refresh(fullPath)
			// 只处理音频文件
			if isMatchFileType(file, fileType) {
				result = append(result, fullPath)
			}
		}
	}
	return result, nil
}

// simple_refresh 极简单行刷新（无依赖，无 Sleep）
func simple_refresh(fullPath string) {
	limit := 80
	displayPath := fullPath
	if len(fullPath) > limit {
		displayPath = "..." + fullPath[len(fullPath)-limit:]
	}
	fmt.Fprintf(os.Stdout, "\r\033[K正在处理：%s", displayPath)
	os.Stdout.Sync()
}

func main() {

	slice := []int{1, 2, 3}

	s1 := slice

	dump.P("%#v", s1)

	return
	list, err := listDir("/Users/weidingyi/workdir", "music", 0)
	fmt.Fprintf(os.Stdout, "\r\033[K")
	if err != nil {
		log.Fatal(err)
	}

	wg := sync.WaitGroup{}
	for _, item := range list {
		wg.Add(1)
		go func(item string) {
			hash, err := getFileInfo(item)
			if err != nil {
				log.Printf("获取文件信息失败: %s, 错误: %v", item, err)
				return
			}
			fmt.Printf("文件: %s 的哈希值为: %s\n", item, hash)
			wg.Done()
		}(item)
	}
	wg.Wait()

}
