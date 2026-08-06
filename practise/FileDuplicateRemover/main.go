package main

import (
	"context"
	"crypto/sha256"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/manifoldco/promptui"
)

const FdrTrashesPath = "./FDRTrashes/"

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
		".mp4a": true, // 苹果设备常用的另一种 M4A 变体
		".opus": true, // 现代网络（如微信、Telegram、语音通话）极常用的高压缩率音频
		".mp2":  true, // 广播电台或老旧 VCD 中常见的音频格式
		".wav_": true, // 某些录音笔或编辑软件产生的临时/备份波形文件
		".mid":  true, // MIDI 音乐文件（虽然体积小，但在老歌库中很常见）
		".midi": true,
		".w64":  true, // 索尼架构下的超大无损音频格式（Sound Foundry 常用）
	},
	"video": {
		".mp4":  true,
		".avi":  true,
		".mov":  true,
		".mkv":  true, // 【核心遗漏】现代高清/4K 电影、多音轨动漫最常用的封装格式
		".flv":  true, // 网页直播、B站早期视频、各种下载器保存下来的流媒体格式
		".f4v":  true, // FLV 的升级版，移动端和部分网页播放器常用
		".rmvb": true, // 【怀旧必备】早年网络下载电影、电视剧的统治级格式
		".rm":   true,
		".wmv":  true, // 微软 Windows Media Video 官方视频格式，老教学片极多
		".webm": true, // 现代网页（如 YouTube）极其推崇的开源高效网页视频格式
		".ts":   true, // 高清电视录制、IPTV 直播流切片格式
		".m3u8": true, // 虽然是 HLS 索引文件，但有时也会被下载器直接存成视频处理
		".mpg":  true, // 经典的 MPEG-1/2 视频，老光盘（VCD/DVD）提取出的视频
		".mpeg": true,
		".vob":  true, // DVD 影碟里的原生视频数据文件
		".3gp":  true, // 早年 2G/3G 时代功能机、老智能机录制视频的格式
		".m4v":  true, // 苹果 iTunes 专用的视频格式
	},
	"image": {
		".jpg":  true, // 最常见的有损压缩格式，几乎所有设备都支持
		".jpeg": true, // JPG 的全称，两者本质相同
		".png":  true, // 支持透明通道的无损格式，网页、截图、设计图最常用
		".gif":  true, // 经典的动图、表情包格式
		".webp": true, // 谷歌推出的现代高效网页图片格式，体积比 JPG 小很多

		".heic": true, // 苹果 iPhone 默认的照片格式（高效率图像容器）
		".heif": true, // HEIC 的标准通用后缀
		".avif": true, // 现代高效图像格式，压缩率比 WebP 更高，新系统和浏览器已全面支持

		".bmp":  true, // Windows 传统的位图格式，无压缩，体积巨大
		".tiff": true, // 印刷和扫描仪常用的高保真无损格式
		".tif":  true, // TIFF 的简写
		".svg":  true, // 矢量图格式，常用于图标、Logo，由于是文本结构，计算哈希可能受格式化影响
		".ico":  true, // Windows 系统图标格式

		// --- 专业摄影 RAW（原始未加工）格式 ---
		".cr2": true, // 佳能 (Canon) 相机
		".cr3": true, // 佳能新一代相机
		".nef": true, // 尼康 (Nikon) 相机
		".arw": true, // 索尼 (Sony) 相机
		".dng": true, // Adobe 统一的通用 RAW 格式（大疆无人机、部分手机专业模式常用）
		".rw2": true, // 松下 (Panasonic) 相机
		".orf": true, // 奥林巴斯 (Olympus) 相机
	},
}

type FileDuplicateRemover struct {
	fileListCh     chan string   // 扇入ch
	workerOutCh    chan FileInfo // 扇出ch
	filePath       string
	fileType       string
	totalScanFiles int64 // 扫描文件计数器
	ctx            context.Context
	cancel         context.CancelFunc
}

type FileInfo struct {
	Name      string
	Size      string
	SizeBites int64
	Hash      string
	ModTime   string
	Path      string
}

func isMatchFileType(file os.DirEntry, fileType string) bool {
	ext := filepath.Ext(file.Name())
	ext = strings.ToLower(ext)
	/**
	ai提示, 这里有潜在panic风险, 如果key不存在就会panic,
	但是实际测试下来, go1.13到go1.25版本都能正常运行, 即便key不存在

	但是绝对不能在一个nil的map中写数据!!!!!!

	官方文档描述:
		"A nil map behaves like an empty map when reading...
		If the map is nil or does not contain such an entry,
		the index expression evaluates to the zero value for the element type of the map.
		"翻译过来就是：一个 nil 的 map 在进行 “读取” 操作时，其行为等同于一个空 map。
		如果这个 map 是 nil 或者找不到对应的 Key，
		索引表达式会安全地返回该 Value 类型的默认零值，
		绝对不会引发错误。
	*/
	return fileTypeMap[fileType][ext]
}

// FormatFileSize 字节转人性化单位 B/KB/MB/GB/TB
func FormatFileSize(bytes int64) string {
	if bytes < 0 {
		return "0 B"
	}
	if bytes < 1024 {
		return fmt.Sprintf("%d B", bytes)
	}

	const unit = 1024
	units := []string{"B", "KB", "MB", "GB", "TB", "PB"}
	val := float64(bytes)
	i := 0

	for val >= unit && i < len(units)-1 {
		val /= unit
		i++
	}

	return fmt.Sprintf("%.2f %s", val, units[i])
}

func MoveFileByCopy(oldPath, newPath string) error {
	_, err := os.Stat(oldPath)
	if errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("%s does not exist", oldPath)
	}

	if err = os.Rename(oldPath, newPath); err == nil { // 尝试同盘移动
		return nil
	}

	/**
	尝试跨盘移动文件
	oldRes, newRes 如果只是依赖两个defer来关闭资源, 那么必须等待程序return之前才能执行defer,
	当运行到os.Remove的时候, windows系统会报待操作的资源未释放,从而Remove失败

	两种解决办法:
		1. 把oldRes, newRes的操作放在匿名自执行函数中, 这样defer必定早于Remove之前运行
		2. 如不使用匿名自执行函数, 则需要在Remove之前, 显式调用close关闭oldRes, newRes
	*/
	copyErr := func() error {
		oldRes, err := os.Open(oldPath)
		if err != nil {
			return err
		}
		defer oldRes.Close()

		// 创建新文件
		newRes, err := os.Create(newPath)
		if err != nil {
			return err
		}
		defer newRes.Close()

		// 使用 io.Copy 进行高效的内核级数据对流（防止大音乐文件卡死内存）
		buf := make([]byte, 64*1024)
		_, err = io.CopyBuffer(newRes, oldRes, buf)
		if err != nil {
			return err
		}
		return nil
	}()

	if copyErr != nil {
		return copyErr
	}

	err = os.Remove(oldPath) // 移除原始文件
	if err != nil {
		return fmt.Errorf("跨盘复制成功，但抹除原文件失败: %w", err)
	}

	return nil

}

func GetFileInfo(filePath string) (FileInfo, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return FileInfo{}, fmt.Errorf("无法打开文件: %s, 错误: %v", filePath, err)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return FileInfo{}, fmt.Errorf("无法获取文件状态: %s, 错误: %v", filePath, err)
	}

	fileSize := fileInfo.Size()
	hasher := sha256.New()

	// 1. 如果文件小于 5MB，直接全量哈希，保证绝对精准
	if fileSize <= 5*1024*1024 {
		if _, err := io.Copy(hasher, file); err != nil {
			return FileInfo{}, fmt.Errorf("读取小文件失败: %s, 错误: %v", filePath, err)
		}
	} else {
		// 2. 如果是大文件，采用多点自适应流式采样
		// 每次只读 64KB，既保证磁头不需要读太多数据，又拓展了采样覆盖面
		const chunkSize = 64 * 1024
		buf := make([]byte, chunkSize)

		// 动态定义 6 个采样比例点：0%(头部), 20%, 40%, 60%, 80%, 100%(尾部减去块大小)
		percentages := []float64{0.0, 0.2, 0.4, 0.6, 0.8}

		// 依次读取这 5 个比例点的数据
		for _, p := range percentages {
			offset := int64(float64(fileSize) * p)
			if _, err := file.ReadAt(buf, offset); err != nil && err != io.EOF {
				return FileInfo{}, err
			}
			hasher.Write(buf)
		}

		// 独立处理最后的尾部块，确保绝对覆盖到文件末尾
		tailOffset := fileSize - chunkSize
		if _, err := file.ReadAt(buf, tailOffset); err != nil && err != io.EOF {
			return FileInfo{}, err
		}
		hasher.Write(buf)

		// 终极加固：注入文件大小。即使内容被伪造，只要大小不同，哈希就绝对不同
		hasher.Write([]byte(fmt.Sprintf("%d", fileSize)))
	}

	f := FileInfo{
		Size:      FormatFileSize(fileSize),
		SizeBites: fileSize,
		Name:      fileInfo.Name(),
		Hash:      fmt.Sprintf("%x", hasher.Sum(nil)),
		ModTime:   fileInfo.ModTime().Format("2006-01-02 15:04:05"),
		Path:      filePath,
	}
	return f, nil
}

func (f *FileDuplicateRemover) walkDir() error {
	/**
	filepath.WalkDir(root, func(path string, d DirEntry, err error))
	root: 待遍历的根路径
	path: 传入这个回调函数内的每个item的绝对路径,即包含root前缀
	d: 传入这个回调函数内的每个item的简要信息
	err: 如果遍历文件出错, 会将这个err传入到这个回调函数内用于判断
	*/
	return filepath.WalkDir(f.filePath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			// 如果是遇到了无权访问的系统目录，优雅跳过该目录，不中断整体扫描
			if d != nil && d.IsDir() {
				log.Printf("🔒 权限不足，已跳过目录: %s (错误: %v)\n", path, err)
				return filepath.SkipDir
			}
			// 如果是个别单文件无法读取，打印后直接跳过该文件
			log.Printf("👻 无法读取文件: %s (错误: %v)\n", path, err)
			return nil
		}
		if !d.IsDir() {
			// 所有的都跳过, 也不应该算作没有扫描到, 所以放在前边
			atomic.AddInt64(&f.totalScanFiles, 1)

			if !isMatchFileType(d, f.fileType) {
				//log.Printf("跳过后缀不合法的文件: %s \n", d.Name())
				return nil // 返回 nil 代表“跳过此错误，继续遍历后面的文件”
			}

			select {
			case <-f.ctx.Done():
				/**
				ai 解释context的cancel()是幂等的, 这里监听到取消信号, 通过f.ctx.Err()返回
				walkDir的调用方, 接受到err != nil 再次执行f.cancel()也是合理合法的
				*/
				return f.ctx.Err()
			case f.fileListCh <- path:
			}
		}
		return nil
	})
}

func NewFileDuplicateRemover(ctx context.Context, filePath, fileType string) *FileDuplicateRemover {
	cancelCtx, cancel := context.WithCancel(ctx) // 根据入参的ctx派生出一个子ctx, 用于控制cancel
	return &FileDuplicateRemover{
		ctx:         cancelCtx,
		cancel:      cancel,
		fileType:    fileType,
		filePath:    filePath,
		fileListCh:  make(chan string, 200),
		workerOutCh: make(chan FileInfo),
	}
}

func (f *FileDuplicateRemover) BuildFileResource() {
	go func() {
		defer close(f.fileListCh)

		defer func() {
			if err := recover(); err != nil {
				log.Printf("捕获到panic: %v\n", err)
				f.cancel() // 终止程序
			}
		}()

		// 递归遍历整个path, 获取所有文件
		err := f.walkDir()
		if err != nil {
			log.Printf("遍历文件错误: %v \n", err)
			f.cancel()
		}
	}()
}

func (f *FileDuplicateRemover) Worker() {
	for path := range f.fileListCh {
		fileInfo, err := GetFileInfo(path) // todo 这里还是同步计算, 如果文件非常大, 导致无法进入监听取消信号的select
		if err != nil {
			log.Printf("获取文件失败: %s, err: %v", path, err)
			continue
		}

		select {
		case <-f.ctx.Done():
			return
		case f.workerOutCh <- fileInfo:
		}
	}
}

func (f *FileDuplicateRemover) preCheck() bool {

	if f.fileType == "" || f.filePath == "" {
		log.Printf("参数有误, fileType: %s, filePath: %s\n", f.fileType, f.filePath)
		return false
	}

	info, err := os.Stat(f.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			log.Printf("❌ 错误：指定的路径不存在: %s\n", f.filePath)
		} else {
			log.Printf("❌ 错误：无法访问指定路径: %s, 错误原因: %v\n", f.filePath, err)
		}
		return false // 返回 false 告知外部终止程序
	}

	// 额外加固：万一用户传进来的路径是个具体的文件（比如 f:/歌曲/1.mp3）而不是文件夹
	if !info.IsDir() {
		log.Printf("❌ 错误：指定的路径不是一个有效的文件夹目录: %s\n", f.filePath)
		return false
	}

	if _, ok := fileTypeMap[f.fileType]; !ok {
		log.Printf("fileType:%s, 不在支持的文件类型容器中\n", f.fileType)
		return false
	}

	return true
}

func (f *FileDuplicateRemover) GoWork() {
	wg := sync.WaitGroup{}
	for i := 0; i < runtime.NumCPU(); i++ {
		/**
		一开始用的这一行构造的workerChs,
		等于从12个长度的基础上又追加了12个长度,
		导致前12个是nil的chan, 导致后续的fanIn读取nil管道panic
		*/
		//workerChs = append(workerChs, w)
		//workerChs[i] = w // 应该用下标精准赋值切片中的每个值

		wg.Add(1)
		go func() {
			defer wg.Done()
			f.Worker()
		}()
	}
	go func() {
		wg.Wait()
		close(f.workerOutCh)
	}()
}

// 返回临时回收站的绝对路径, err
func (f *FileDuplicateRemover) initTrashesDir() (string, error) {
	// 使用局部变量，避免重复调用时不断拼接 f.trashesPath 导致路径错误叠加
	trashPath := filepath.Join(FdrTrashesPath, f.fileType, time.Now().Format("20060102"))

	if err := os.MkdirAll(trashPath, 0777); err != nil { // If path is already a directory, MkdirAll does nothing
		return "", err
	}

	return filepath.Abs(trashPath)
}

func (f *FileDuplicateRemover) OperateDuplicateFile() {
	fileMap := make(map[string][]FileInfo)
	for item := range f.workerOutCh {
		cleanHash := strings.ToLower(strings.TrimSpace(item.Hash)) // 去除hash两端的空字符
		fileMap[cleanHash] = append(fileMap[cleanHash], item)
	}

	duplicateMap := make(map[string][]FileInfo)
	for key, val := range fileMap {
		if len(val) > 1 {
			duplicateMap[key] = val
		}
	}

	totalScanned := atomic.LoadInt64(&f.totalScanFiles)

	msg := fmt.Sprintf("👀 扫描完成！共扫描了 %d 个文件, ", totalScanned)

	if len(duplicateMap) > 0 {
		msg += fmt.Sprintf("发现 %d 组重复文件", len(duplicateMap))
		log.Println(msg)
	} else {
		msg += fmt.Sprintln("没有重复文件")
		log.Println(msg)
		return
	}

	// 在本地创建临时回收站目录
	absTrashPath, err := f.initTrashesDir()
	if err != nil {
		log.Println(err)
		return
	}

	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "\U0001F336 {{ .Name | cyan }}",
		Inactive: "  {{ .Name | cyan }}",
		Selected: "\U0001F336 已删除: {{ .Path | red | cyan }}",
		Details: `
------------------ 详细信息 ---------------------
{{ "Path:" | faint }}	{{ .Path }}
{{ "Size:" | faint }}	{{ .Size }}
{{ "ModTime:" | faint }}	{{ .ModTime }}
`,
	}
	selectedCnt := 0
	wg := sync.WaitGroup{}

	prompt := promptui.Select{
		Label:     fmt.Sprintf("%s 选择一个待删除的文件并按下回车键", time.Now().Format("2006-01-02 15:04:05")),
		Templates: templates,
		HideHelp:  true,
	}

	for _, infos := range duplicateMap {
		prompt.Items = infos

		selectedIdx, _, err := prompt.Run()
		if err != nil {
			fmt.Printf("退出程序: %v\n", err)
			return
		}

		selectedFileInfo := infos[selectedIdx]

		// 检查文件是否还存在, 存在则移动至当前路径下的回收站,
		wg.Add(1)
		go func(selectedFile FileInfo) {
			defer wg.Done()

			err = MoveFileByCopy(selectedFile.Path, filepath.Join(absTrashPath, selectedFile.Name))
			if err != nil {
				log.Printf("MoveFileByCopy: %v\n", err)
			}

			// todo 构造一键还原文件的结构体容器, 生成reset.json在trash目录下, 留待一键还原文件

		}(selectedFileInfo)

		selectedCnt += 1
	}

	wg.Wait() // 等待所有文件都移动完毕

	log.Printf("已删除%d个文件, 原始文件备份至: %s\n", selectedCnt, absTrashPath)
	return
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("panic: %v\n", err)
			return
		}
	}()

	// 1. 定义命令行参数变量
	var pathParam string
	var typeParam string

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "   File Duplicate Remover (fdr) - 重复文件清理工具\n")
		fmt.Fprintf(os.Stderr, "=========================================================\n")
		fmt.Fprintf(os.Stderr, "描述:\n")
		fmt.Fprintf(os.Stderr, "  本工具用于扫描指定目录下的文件，通过计算 SHA256哈希值。\n")
		fmt.Fprintf(os.Stderr, "  精准定位重复文件，并通过交互式命令行界面安全地进行清理。\n")
		fmt.Fprintf(os.Stderr, "  清理完自动保留备份至%s, 保留手动操作空间。\n\n", FdrTrashesPath)
		fmt.Fprintf(os.Stderr, "用法:\n")
		fmt.Fprintf(os.Stderr, "  %s -p <目录路径> -t <文件类型>\n\n", filepath.Base(os.Args[0]))
		fmt.Fprintf(os.Stderr, "参数说明:\n")
		flag.PrintDefaults() // 自动打印下方定义的参数和默认值
		fmt.Fprintf(os.Stderr, "\n")
	}
	// 3. 绑定命令行参数
	flag.StringVar(&pathParam, "p", "", "✨【必填】需要扫描的目标文件夹路径 (例如: F:\\歌曲)")
	flag.StringVar(&typeParam, "t", "", "✨【必填】扫描的文件类型: music(音频), video(视频)\nimage(图片), 程序"+
		"内置一批常见文件类型后缀")

	// 4. 解析参数
	flag.Parse()

	// 5. 校验必填参数
	if pathParam == "" {
		log.Printf("👻 请使用 `-p` 指定扫描的路径\n")
		return
	}
	if typeParam == "" {
		log.Printf("👻 请使用 `-t` 指定扫描的类型\n")
		return
	}

	ctx := context.Background()

	fdr := NewFileDuplicateRemover(ctx, pathParam, typeParam)

	// 确保main函数执行结束后关闭各个资源
	defer fdr.cancel()

	// 预检查
	if !fdr.preCheck() {
		return
	}

	// 构造文件资源通道
	fdr.BuildFileResource()

	// 扇出: 启动N个worker并发计算文件哈希
	fdr.GoWork()

	// 扇入: 交互式操作重复文件
	fdr.OperateDuplicateFile()
	return
}
