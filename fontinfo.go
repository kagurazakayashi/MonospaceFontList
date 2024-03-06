package main

import (
	"fmt"
	"log"
	"os"

	"encoding/json"

	"golang.org/x/image/font/sfnt"
)

func fontInfo() {
	for i, fontinfo := range fontPathList {
		var ii string = fmt.Sprintf("%d / %d", i+1, fontPathListLen)
		fontFile, err := os.ReadFile(fontinfo.FontPath)

		fmt.Printf("\r\033[K正在获取字体信息: %s", ii)

		if err != nil {
			log.Printf("错误：无法打开文件 %s : %v", fontinfo.FontPath, err)
		}

		font, err := sfnt.Parse(fontFile)
		if err != nil {
			log.Printf("%s 错误：解析字体文件 %s 失败: %v", ii, fontinfo.FontPath, err)
			continue
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
		fontPathList[i] = fontinfo
	}
}

func outJSONInfo() {
	jsonData, err := json.MarshalIndent(fontPathList, "", "  ")
	if err != nil {
		fmt.Println("创建 JSON 失败:", err)
		return
	}
	log.Println(string(jsonData))
}
