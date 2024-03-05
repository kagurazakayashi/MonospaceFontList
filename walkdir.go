package main

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func fileList() {
	err := filepath.Walk(scanDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			ext := strings.ToLower(filepath.Ext(info.Name()))
			if ext == ".ttf" || ext == ".otf" {
				absPath, err := filepath.Abs(path)
				if err != nil {
					log.Println("错误：无法访问文件夹:", err)
					return err
				}
				fontPathList = append(fontPathList, []string{absPath})
			}
		}
		return nil
	})
	if err != nil {
		log.Println("错误：无法遍历文件夹:", err)
	} else {
		log.Println("搜索到字体:", len(fontPathList))
	}
}

func systemFontDir() string {
	switch runtime.GOOS {
	case "windows":
		return os.Getenv("WINDIR") + "\\Fonts"
	case "darwin":
		return "/Library/Fonts"
	case "linux":
		return "/usr/share/fonts"
	default:
		log.Println("错误：不支持的操作系统。")
		os.Exit(1)
		return ""
	}
}
