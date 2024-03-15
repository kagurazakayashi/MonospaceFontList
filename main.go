package main

import (
	"flag"
	"log"
	"os"
	"time"
)

type FontInfo struct {
	FontPath      string
	Name          string
	Family        string
	SubfamilyName string
	Version       string
	UnitsPerEm    int
	Monospaced    int // -1: 不等宽, 0: 未知, 1+: 等宽宽度
	MonospacedZH  int // -2: 不支持中文, -1: 不等宽, 0: 未知, 1+: 等宽宽度
	MD5           string
}

var (
	scanDir          string
	extensionNames   string
	timeLayout       string     = "2006-01-02 15:04:05"
	fontPathList     []FontInfo = []FontInfo{}
	fontPathListLen  int
	enTestChars      string
	zhTestChars      string
	chineseTotal     uint = 0
	chineseMonoTotal uint = 0
	chineseChkLevel  int
	monospacedTotal  uint = 0
	fontSize         float64
	outDir           string
	imageTop         int
	imageLeft        int
	imageWidth       int
	imageHeight      int
	imageText        string
	pageTitle        string
)

func main() {
	log.Println("MonospaceFontList 0.0.1  " + time.Now().Format(timeLayout))
	loadHTML()

	flag.StringVar(&scanDir, "i", "", "要扫描的字体文件夹，默认为系统字体文件夹。")
	flag.StringVar(&extensionNames, "e", "ttf,otf,ttc", "要扫描的字体文件扩展名，用 `,` 分隔。默认为 `ttf,otf,ttc` 。")
	flag.StringVar(&enTestChars, "en", "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789", "英文测试字符列表，默认大小写字母和数字。")
	flag.StringVar(&zhTestChars, "zt", "正，一丨中！小。", "中文测试字符列表。")
	flag.IntVar(&chineseChkLevel, "zh", 7, "中文 Unicode 区块检查级别( 0.基本 - 7.扩展 )")
	flag.StringVar(&outDir, "o", "out", "输出文件夹，不带结尾 `/` 。默认为 `当前文件夹/out` 。")
	flag.Float64Var(&fontSize, "s", 64, "字体预览图所使用的字体大小")
	flag.IntVar(&imageLeft, "x", 10, "字体预览图左侧空白")
	flag.IntVar(&imageTop, "y", 64, "字体预览图顶部空白")
	flag.IntVar(&imageWidth, "w", 800, "字体预览图宽度")
	flag.IntVar(&imageHeight, "h", 78, "字体预览图高度")
	flag.StringVar(&imageText, "t", "AaBbCc0123?.文字。", "字体预览图文字")
	flag.StringVar(&pageTitle, "p", "MonospaceFontList", "网页标题前缀")
	flag.Parse()

	if len(scanDir) == 0 {
		scanDir = systemFontDir()
	}
	if len(extensionNames) == 0 {
		log.Println("错误：你必须指定一个扩展名。")
		os.Exit(1)
	}
	fileList()
	fontInfo()
	saveHTML()
	// outJSONInfo()
	log.Println("字体信息收集完成。")
	log.Println("字体数量:", fontPathListLen, "中文字体数量:", chineseTotal)
	log.Println("等宽字体数量:", monospacedTotal, "中文等宽字体数量:", chineseMonoTotal)
}
