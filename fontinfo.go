package main

import (
	"crypto/md5"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"os"
	"sync"

	"encoding/json"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/font/sfnt"
	"golang.org/x/image/math/fixed"
)

var (
	mu sync.Mutex
	wg sync.WaitGroup
)

func fontInfo() {
	if maxGoroutines == 1 {
		for i := range fontPathList {
			fontInfoN, tpp := work(i)
			workSave(i, fontInfoN, tpp)
			fmt.Printf("\r\033[K已完成: %d / %d (%d %%)  ", total, fontPathListLen, total*100/uint(fontPathListLen))
		}
	} else {
		var goroutineSem chan struct{} = make(chan struct{}, maxGoroutines)
		for i := range fontPathList {
			goroutineSem <- struct{}{}
			wg.Add(1)
			go func() {
				defer wg.Done()
				fontInfoN, tpp := work(i)
				mu.Lock()
				workSave(i, fontInfoN, tpp)
				fmt.Printf("\r\033[K已完成: %d / %d (%d %%) [并行 %d]  ", total, fontPathListLen, total*100/uint(fontPathListLen), len(goroutineSem))
				mu.Unlock()
				<-goroutineSem
			}()
		}
		wg.Wait()
	}
}

func workSave(i int, fontInfoN FontInfo, tpp [4]uint) {
	fontPathList[i] = fontInfoN
	total += tpp[0]
	monospacedTotal += tpp[1]
	chineseTotal += tpp[2]
	chineseMonoTotal += tpp[3]
}

func work(i int) (FontInfo, [4]uint) {
	var fontinfo FontInfo = fontPathList[i]

	// 0okTotal,1monospacedTotal,2chineseTotal,3chineseMonoTotal
	var tpp [4]uint = [4]uint{0, 0, 0, 0}

	var ii string = fmt.Sprintf("%d / %d", i+1, fontPathListLen)
	fontFile, err := os.ReadFile(fontinfo.FontPath)

	if err != nil {
		log.Printf("错误：无法打开文件 %s : %v", fontinfo.FontPath, err)
		tpp[0]++
		return fontinfo, tpp
	}

	font, err := sfnt.Parse(fontFile)
	if err != nil {
		log.Printf("%s 错误：解析字体文件 %s 失败: %v", ii, fontinfo.FontPath, err)
		tpp[0]++
		return fontinfo, tpp
	}

	fontinfo.Name, err = font.Name(nil, sfnt.NameIDFull)
	if err != nil {
		log.Printf("%s 错误：获取字体名称 %s 失败: %v", ii, fontinfo.FontPath, err)
	}

	fontinfo.Family, err = font.Name(nil, sfnt.NameIDFamily)
	if err != nil {
		log.Printf("%s 错误：获取字体家族 %s 失败: %v", ii, fontinfo.FontPath, err)
	}

	fontinfo.SubfamilyName, err = font.Name(nil, sfnt.NameIDSubfamily)
	if err != nil {
		log.Printf("%s 错误：获取字体样式 %s 失败: %v", ii, fontinfo.FontPath, err)
	}

	fontinfo.Version, err = font.Name(nil, sfnt.NameIDVersion)
	if err != nil {
		log.Printf("%s 错误：获取字体版本 %s 失败: %v", ii, fontinfo.FontPath, err)
	}

	fontinfo.UnitsPerEm = int(font.UnitsPerEm())
	fontinfo.Monospaced = isMonospaced(font, enTestChars)
	if isSupportsChinese(font) {
		tpp[2]++
		fontinfo.MonospacedZH = isMonospaced(font, zhTestChars)
		if fontinfo.Monospaced > 0 && fontinfo.MonospacedZH > 0 {
			tpp[1]++
			tpp[3]++
		}
	} else {
		fontinfo.MonospacedZH = -2
		if fontinfo.Monospaced > 0 {
			tpp[1]++
		}
	}

	fontinfo.MD5 = fmt.Sprintf("%x", md5.Sum(fontFile))
	drew(fontFile, fontinfo.MD5)
	tpp[0]++
	return fontinfo, tpp
}

func outJSONInfo() {
	jsonData, err := json.MarshalIndent(fontPathList, "", "  ")
	if err != nil {
		fmt.Println("创建 JSON 失败:", err)
		return
	}
	log.Println(string(jsonData))
}

func isSupportsChinese(font *sfnt.Font) bool {
	// 定义中文Unicode区块
	var chineseUnicodeRanges [][]int = [][]int{
		{0x4E00, 0x9FFF},   // 0基本汉字
		{0x3400, 0x4DBF},   // 1扩展A
		{0x20000, 0x2A6DF}, // 2扩展B
		{0x2A700, 0x2B73F}, // 3扩展C
		{0x2B740, 0x2B81F}, // 4扩展D
		{0x2B820, 0x2CEAF}, // 5扩展E
		{0x2CEB0, 0x2EBEF}, // 6扩展F
		{0x30000, 0x3134F}, // 7扩展G
	}
	var isOK int = 0
	for i, rangeValue := range chineseUnicodeRanges {
		if i > chineseChkLevel {
			break
		}
		for runeValue := rune(rangeValue[0]); runeValue <= rune(rangeValue[1]); runeValue++ {
			r, err := font.GlyphIndex(&sfnt.Buffer{}, runeValue)
			if err == nil && r != 0 {
				isOK++
				break
			}
		}
	}
	return isOK == chineseChkLevel+1
}

func isMonospaced(font *sfnt.Font, testChars string) int {
	var oldWidth fixed.Int26_6 = 0
	for i, ch := range testChars {
		buf := &sfnt.Buffer{}
		mIndex, err := font.GlyphIndex(buf, ch)
		if err != nil {
			log.Printf("错误：获取字符“ %c ”的索引失败: %v", ch, err)
		}
		mWidth, err := font.GlyphAdvance(buf, mIndex, fixed.Int26_6(font.UnitsPerEm()), 0)
		if err != nil {
			log.Printf("错误：获取字符“ %c ”的宽度失败: %v", ch, err)
		}
		if i == 0 {
			oldWidth = mWidth
		} else if oldWidth != mWidth {
			return -1
		}
	}
	return int(oldWidth)
}

func drew(fontFile []byte, name string) {
	f, err := opentype.Parse(fontFile)
	if err != nil {
		log.Println("解析字体失败:", name, err)
		return
	}

	face, err := opentype.NewFace(f, &opentype.FaceOptions{
		Size:    fontSize,
		DPI:     fontDPI,
		Hinting: font.Hinting(fontHinting),
	})
	if err != nil {
		log.Println("创建字体face失败:", name, err)
		return
	}
	defer face.Close()

	img := image.NewRGBA(image.Rect(0, 0, imageWidth, imageHeight))
	draw.Draw(img, img.Bounds(), image.NewUniform(color.White), image.Point{}, draw.Src)

	point := fixed.Point26_6{
		X: fixed.I(imageLeft),
		Y: fixed.I(imageTop),
	}

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(color.Black),
		Face: face,
		Dot:  point,
	}
	d.DrawString(imageText)

	if _, err := os.Stat(outDir); os.IsNotExist(err) {
		os.Mkdir(outDir, os.ModePerm)
	}
	var pngFile string = outDir + "/" + name + ".png"
	// pngFile = strings.Replace(pngFile, " ", "_", -1)
	outFile, err := os.Create(pngFile)
	if err != nil {
		log.Println("创建图片文件失败:", err)
		return
	}
	defer outFile.Close()

	err = png.Encode(outFile, img)
	if err != nil {
		log.Println("图片保存失败:", err)
		return
	}
	// log.Println("图片已保存为 ", pngFile)
}
