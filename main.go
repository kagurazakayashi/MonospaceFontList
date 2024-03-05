package main

import (
	"flag"
	"log"
	"os"
	"time"
)

var (
	scanDir        string
	extensionNames string
	timeLayout     string     = "2006-01-02 15:04:05"
	fontPathList   [][]string = [][]string{}
)

func main() {
	log.Println("MonospaceFontList 0.0.1  " + time.Now().Format(timeLayout))
	flag.StringVar(&scanDir, "i", "", "要扫描的字体文件夹，默认为系统字体文件夹")
	flag.StringVar(&extensionNames, "e", "ttf,otf,ttc", "要扫描的字体文件扩展名，用 `,` 分隔。默认为 `ttf,otf,ttc` 。")
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
}
