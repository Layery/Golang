package image

import (
	"image"
	"image/jpeg"
	"os"
	"strings"
	"wisdom/controller"

	"github.com/fogleman/gg"
	"github.com/gin-gonic/gin"
)

type Image struct {
	controller.Base
}

// 排版配置：固定 3:4 竖图 1080 × 1440
type LayoutConfig struct {
	CanvasWidth  float64
	CanvasHeight float64
	MarginLeft   float64
	MarginRight  float64
	MarginTop    float64
	MarginBottom float64
	CardRadius   float64
	LineHeight   float64 // 行间距
	SectionGap   float64 // 段间距
}

func NewImage() *Image {
	return &Image{}
}

func (img *Image) GenImage(c *gin.Context) {
	// 1. 打开无文字底图
	inputPath := "./assets/images/base_map_no_text.jpg"
	outputPath := "poster_with_text.jpg"

	file, err := os.Open(inputPath)
	if err != nil {
		panic("打开底图失败: " + err.Error())
	}
	defer file.Close()

	baseImg, _, err := image.Decode(file)
	if err != nil {
		panic("解码底图失败: " + err.Error())
	}

	// 2. 创建绘图上下文
	dc := gg.NewContextForImage(baseImg)

	// 3. 排版参数
	cfg := LayoutConfig{
		CanvasWidth:  float64(dc.Width()),
		CanvasHeight: float64(dc.Height()),
		MarginLeft:   90,
		MarginRight:  90,
		MarginTop:    420,
		MarginBottom: 40,
		CardRadius:   24,
		LineHeight:   100,
		SectionGap:   80,
	}

	// 4. 正文白色半透明卡片
	cardX := cfg.MarginLeft
	cardY := cfg.MarginTop
	cardW := cfg.CanvasWidth - cfg.MarginLeft - cfg.MarginRight
	cardH := cfg.CanvasHeight - cfg.MarginTop - cfg.MarginBottom

	dc.SetRGBA(1, 1, 1, 100)
	dc.DrawRoundedRectangle(cardX, cardY, cardW, cardH, cfg.CardRadius)
	dc.Fill()

	// 5. 标题
	titleFont := XingKaiFont
	err = dc.LoadFontFace(titleFont, 300)
	if err != nil {
		panic("加载标题字体失败: " + err.Error())
	}

	dc.SetRGB(0.32, 0.18, 0.12)
	//title := "藏在生活中的智慧"
	title := "老人言"
	dc.DrawStringAnchored(title, cfg.CanvasWidth/2, 230, 0.5, 0.5)

	// 6. 副标题
	//err = dc.LoadFontFace(YaHeiFont, 80)
	//if err != nil {
	//	panic("加载副标题字体失败: " + err.Error())
	//}
	//dc.SetRGB(0.4, 0.25, 0.2)
	//subtitle := "人活到五六十，越品越觉得老话靠谱"
	//dc.DrawStringAnchored(subtitle, cfg.CanvasWidth/2, 580, 0.5, 0.5)

	// 7. 正文
	err = dc.LoadFontFace(XinWeiFont, 100)
	if err != nil {
		panic("加载正文字体失败: " + err.Error())
	}

	dc.SetRGB(0.2, 0.2, 0.2)

	bodyText := `
一个人真正的觉醒，
是对所有关系的绝望，
当你见过人性的黑暗，
遭遇过信任的崩塌，
体验过孤立无援的滋味，
你就会明白人生的本质
就是一个人活着，
不要高估感情，
不要低估人性。
这个世界最没意思的就是人了。
`
	textStartX := cardX + 20
	textStartY := cardY + 10
	maxTextWidth := cardW - 40
	currentY := textStartY

	paragraphs := strings.Split(bodyText, "\n")

	for _, para := range paragraphs {
		if strings.TrimSpace(para) == "" {
			currentY += cfg.SectionGap
			continue
		}

		dc.DrawStringWrapped(
			para,
			textStartX,
			currentY,
			0, 0,
			maxTextWidth,
			cfg.LineHeight,
			gg.AlignCenter,
		)
		currentY += cfg.LineHeight + cfg.SectionGap
	}

	// 8. 底部标签
	tags := []string{"生活感悟", "平和心态"}
	tags = []string{}
	tagWidth := float64(160)
	tagHeight := float64(36)
	tagGap := float64(20)
	totalWidth := len(tags)*int(tagWidth) + (len(tags)-1)*int(tagGap)
	startX := (int(cfg.CanvasWidth) - totalWidth) / 2
	tagY := cfg.CanvasHeight - 60

	for i, tag := range tags {
		tx := float64(startX + i*int((tagWidth+tagGap)))

		dc.SetRGB(0.76, 0.6, 0.42)
		dc.DrawRoundedRectangle(tx, tagY, tagWidth, tagHeight, 18)
		dc.Fill()

		dc.SetRGB(1, 1, 1)
		dc.DrawStringAnchored(tag, tx+tagWidth/2, tagY+tagHeight/2, 0.5, 0.5)
	}

	// 9. 保存图片
	outFile, err := os.Create(outputPath)
	if err != nil {
		panic("创建输出文件失败: " + err.Error())
	}
	defer outFile.Close()

	err = jpeg.Encode(outFile, dc.Image(), &jpeg.Options{Quality: 90})
	if err != nil {
		panic("保存图片失败: " + err.Error())
	}

	println("生成成功:", outputPath)
}
