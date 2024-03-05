package main

import (
	"log"
	"os"
)

func fontInfo() {
	fontPath := "1.ttf"
	fontFile, err := os.Open(fontPath)
	if err != nil {
		log.Fatalf("错误：无法打开文件: %v", err)
	}
	defer fontFile.Close()
}
