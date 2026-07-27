package image

import (
	"fmt"
	"image/jpeg"
	"os"
	"strings"
	"wisdom/controller"

	"github.com/fogleman/gg"
	"github.com/gin-gonic/gin"
)

type CommonImage struct {
	controller.Base
}

func NewCommonImage() *CommonImage {
	image := new(CommonImage)
	return image
}

// 通用文案载体，全模板通用
type PosterData struct {
	Title        string     // 顶部大标题
	HighlightTxt string     // 高亮文字（无则空）
	Paragraphs   []string   // 单栏正文段落
	ListItems    []string   // 清单条目（打分图专用）
	LeftColList  []LineItem // 双栏左列文案（对比图专用）
	RightColList []LineItem // 双栏右列文案
	Tags         []string   // 底部标签
	DecorPNGPath string     // 装饰插画路径（无则空）
}

type LineItem struct {
	Age string
	Txt string
}

type TemplateConfig struct {
	// 画布基础
	W, H     float64
	BgMode   string // pure纯色 / grid网格 / image底图
	BgColor  []float64
	GridStep float64 // 网格间距，grid模式生效
	BgImg    string

	// 全局边距、文字基础
	MarginL, MarginR, MarginT, MarginB float64
	TextPaddingX                       float64
	LineH, SectionGap                  float64

	// 标题配置
	ShowTitle    bool
	TitleVisH    float64
	TitleFontKey string
	TitleColor   []float64
	TitleStroke  float64
	TitleAlign   string // center/left

	// 卡片
	ShowCard     bool
	CardRGBA     []float64
	CardRadius   float64
	CardMaxTextH float64 // 文字最大高度，防溢出
	AutoCutText  bool    // 超长自动截断

	// 布局模式（核心，区分所有版式）
	LayoutType string
	// 可选值：
	// single_card 单栏古风卡片
	// double_col 左右双栏对比
	// score_list 清单打分列表

	// 装饰元素
	DecorOffsetX, DecorOffsetY float64
	DecorW, DecorH             float64
}

// 唯一对外入口，传入配置+文案，自动生成任意版式海报
/*func RenderPoster(cfg TemplateConfig, data PosterData) error {
	// ========== 公共逻辑：所有版式都执行，只写1次 ==========
	// 1. 创建画布、绘制背景（纯色/网格/底图统一封装）
	dc, err := NewDrawCanvas(cfg)
	if err != nil {
		return err
	}
	// 2. 预加载字体库，自适应字号换算
	fontSet, err := LoadFontSet(cfg)
	if err != nil {
		return err
	}
	// 3. 绘制白色卡片（开关控制，所有版式通用）
	if cfg.ShowCard {
		DrawCard(dc, cfg)
	}
	// ========== 布局分支：仅几十行区分逻辑，复用公共绘图 ==========
	switch cfg.LayoutType {
	case "single_card":
		// 古风单栏卡片：调用通用安全文字绘制
		DrawSingleColumnText(dc, fontSet.Body, cfg, data.Paragraphs)
	case "double_col":
		// 左右对比分栏：调用通用双栏绘制工具
		DrawDoubleCompareCol(dc, fontSet.Body, cfg, data.LeftColList, data.RightColList)
	case "score_list":
		// 打分清单版式：循环绘制列表条目
		DrawScoreList(dc, fontSet.Body, cfg, data.ListItems)
	}
	// ========== 公共后置元素：所有版式通用 ==========
	// 绘制顶部标题
	if cfg.ShowTitle {
		DrawStrokeTitle(dc, fontSet.Title, cfg, data.Title)
	}
	// 叠加右下角装饰插画
	if data.DecorPNGPath != "" {
		DrawDecorPNG(dc, cfg, data.DecorPNGPath)
	}
	// 底部标签、高亮文字、保存图片
	DrawBottomTags(dc, cfg, data.Tags)
	return SavePoster(dc, "out.jpg")
}
*/
// WrapAutoLine 输入绘图上下文、文本、单行最大安全宽度，返回自动拆分后的多行
// 自动适配当前加载的任意中文字体，不会横向溢出卡片
func WrapAutoLine(dc *gg.Context, text string, maxSafeW float64) []string {
	var resultLines []string
	var currentRunes []rune

	for _, r := range []rune(text) {
		// 测试追加当前字符后的总宽度
		tempStr := string(append(currentRunes, r))
		strW, _ := dc.MeasureString(tempStr)
		// 超过安全宽度，换行
		if strW > maxSafeW && len(currentRunes) > 0 {
			resultLines = append(resultLines, string(currentRunes))
			currentRunes = []rune{r}
		} else {
			currentRunes = append(currentRunes, r)
		}
	}
	// 剩余文字收尾
	if len(currentRunes) > 0 {
		resultLines = append(resultLines, string(currentRunes))
	}
	return resultLines
}

// CalcTextTotalHeight 计算所有段落渲染需要的总纵向高度
func CalcTextTotalHeight(dc *gg.Context, paragraphs []string, maxSafeW, lineH, segGap float64) float64 {
	totalH := 0.0
	for _, p := range paragraphs {
		p = strings.TrimSpace(p)
		if p == "" {
			totalH += segGap
			continue
		}
		lines := WrapAutoLine(dc, p, maxSafeW)
		// 本段高度 = 行数 × 单行高度
		totalH += float64(len(lines)) * lineH
		// 段落间距
		totalH += segGap
	}
	return totalH
}

// SafeDrawTextInCard 在卡片安全区域内绘制文字，自动换行、超长截断
// 返回两个值：渲染结束的Y坐标、是否发生文字截断
func SafeDrawTextInCard(
	dc *gg.Context,
	paragraphs []string,
	safeX, startY, maxSafeW, maxSafeH float64,
	lineH, segGap float64,
	textColor []float64,
) (float64, bool) {
	dc.SetRGB(textColor[0], textColor[1], textColor[2])
	currentY := startY
	remainH := maxSafeH // 卡片剩余可绘制高度
	isTruncated := false

	for _, para := range paragraphs {
		paraTrim := strings.TrimSpace(para)
		if paraTrim == "" {
			if remainH < segGap {
				isTruncated = true
				break
			}
			currentY += segGap
			remainH -= segGap
			continue
		}
		// 根据当前字体自动分行，横向锁死maxSafeW
		lineList := WrapAutoLine(dc, paraTrim, maxSafeW)
		for _, line := range lineList {
			// 剩余高度不足以绘制一行，直接截断
			if remainH < lineH {
				dc.DrawString("******", safeX, currentY)
				isTruncated = true
				return currentY + lineH, isTruncated
			}
			// 绘制本行，绝对不会超出左右安全宽度
			dc.DrawString(line, safeX, currentY)
			currentY += lineH
			remainH -= lineH
		}
		// 段落间距留白
		if remainH < segGap {
			isTruncated = true
			break
		}
		currentY += segGap
		remainH -= segGap
	}
	return currentY, isTruncated
}

func (image *CommonImage) GenCommonImage(ctx *gin.Context) {

	// 画布尺寸 960 × 1283 和示例打分图一致
	const w, h = 960.0, 1283.0
	fontPath := XinWeiFont

	// 1. 创建画布，绘制浅紫色网格背景
	dc := gg.NewContext(int(w), int(h))
	// 底色
	dc.SetRGB(0.92, 0.86, 0.96)
	dc.Clear()
	// 网格线
	dc.SetRGBA(1, 1, 1, 180)
	gridStep := 30.0
	for x := 0.0; x < w; x += gridStep {
		dc.DrawLine(x, 0, x, h)
		dc.SetLineWidth(1)
		dc.Stroke()
	}
	for y := 0.0; y < h; y += gridStep {
		dc.DrawLine(0, y, w, y)
		dc.SetLineWidth(1)
		dc.Stroke()
	}

	// 2. 卡片基础参数
	margin := 60.0
	cardX := margin
	cardY := 60.0
	cardW := w - margin*2
	cardH := h - margin*2
	// 文字距离卡片内边距（安全边界）
	textPaddingX := 40.0
	textPaddingY := 40.0
	// 文字安全绘制范围（核心约束，文字不会超出卡片）
	textSafeX := cardX + textPaddingX
	textSafeMaxW := cardW - textPaddingX*2
	textSafeStartY := cardY + textPaddingY
	textSafeMaxH := cardH - textPaddingY*2

	// 绘制白色圆角卡片
	dc.SetRGBA(1, 1, 1, 245)
	dc.DrawRoundedRectangle(cardX, cardY, cardW, cardH, 24)
	dc.Fill()

	// 3. 绘制顶部大标题
	err := dc.LoadFontFace(fontPath, 56)
	if err != nil {
		panic(fmt.Sprintf("加载标题字体失败: %v", err))
	}
	dc.SetRGB(0.55, 0.25, 0.75)
	titleText := "能达到 60 分，就千万别离婚……"
	// 标题居中自动换行
	dc.DrawStringWrapped(titleText, cardX+40, cardY+50, 0, 0, cardW-80, 70, gg.AlignLeft)

	// 4. 加载正文字体，设置基础排版参数
	err = dc.LoadFontFace(fontPath, 32)
	if err != nil {
		panic(fmt.Sprintf("加载正文字体失败: %v", err))
	}
	lineHeight := 50.0
	segGap := 10.0
	textColor := []float64{0.2, 0.2, 0.2}

	// 模拟两套文案，切换测试长短适配
	// 短文案（完整展示无截断）
	//listParagraphs := []string{
	//	"不家暴（10分）",
	//	"能养家（10分）",
	//	"不出轨（10分）",
	//}
	// 超长文案（会自动截断并提示）
	listParagraphs := []string{
		"不家暴（10分）",
		"能养家（10分）",
		"不出轨（10分）",
		"愿意做家务（10分）",
		"不乱花钱（10分）",
		"尊重你的父母（10分）",
		"不嫌弃你（10分）",
		"有话好好说（10分）",
		"周末愿意陪你（10分）",
		"记得你的生日记得你的生日记得你的生日记得你的生日记得你的生日记得你的生日记得你的生日记得你的生日记得你的生日记得你的生日记得你的生日记得你的生日记得你的生日记得你的生日记得你的生日记得你的生日记得你的生日记得你的生日记得你的生日（10分）",
		"愿意做家务（10分）",
		"不乱花钱（10分）",
		"尊重你的父母（10分）",
		"不嫌弃你（10分）",
		"有话好好说（10分）",
		"周末愿意陪你（10分）",
		"记得你的生日（10分）",
		"遇事主动沟通，不冷暴力（10分）",
		"愿意分担家庭开销，不自私（10分）",
		"懂得包容你的缺点，不随意指责（10分）",
	}

	// 工具2：预计算文字总高度（可选，用于动态卡片高度）
	totalTextH := CalcTextTotalHeight(dc, listParagraphs, textSafeMaxW, lineHeight, segGap)
	fmt.Printf("文字总占用高度：%.2f px\n", totalTextH)

	// 工具3：安全绘制文字，锁死卡片边界，超长自动截断
	endY, truncated := SafeDrawTextInCard(
		dc,
		listParagraphs,
		textSafeX,
		textSafeStartY,
		textSafeMaxW,
		textSafeMaxH,
		lineHeight,
		segGap,
		textColor,
	)

	// 打印适配结果
	if truncated {
		fmt.Println("⚠️ 警告：文案超出卡片纵向高度，已自动截断末尾文字")
	} else {
		fmt.Println("✅ 文案完整展示，无溢出")
	}
	fmt.Printf("文字绘制结束Y坐标：%.2f\n", endY)

	// 右上角绘制小皇冠装饰
	dc.SetRGB(0.95, 0.65, 0.22)
	dc.LoadFontFace(fontPath, 40)
	dc.DrawString("♛", cardX+cardW-120, cardY+90)

	// 导出图片
	outFile, err := os.Create("score_list_output.jpg")
	if err != nil {
		panic(fmt.Sprintf("创建输出文件失败: %v", err))
	}
	defer outFile.Close()
	err = jpeg.Encode(outFile, dc.Image(), &jpeg.Options{Quality: 95})
	if err != nil {
		panic(fmt.Sprintf("保存图片失败: %v", err))
	}
	fmt.Println("图片生成完成：score_list_output.jpg")
}
