package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

var (
	htmlTemp     string = ""
	htmlBodyTemp string = ""
)

func loadHTML() {
	var htmlPath string = "html/list.html"
	file, err := os.Open(htmlPath)
	if err != nil {
		log.Println("错误：无法读取资源文件: ", err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var mode int8 = 1
	for scanner.Scan() {
		var line string = scanner.Text()
		if mode == 1 && strings.Contains(line, "<!-- item -->") {
			mode = 20
		} else if mode == 2 && strings.Contains(line, "<!-- !item -->") {
			mode = 10
		}
		switch mode {
		case 10:
			mode = 1
		case 1:
			htmlTemp = htmlTemp + "\n" + line
		case 20:
			mode = 2
		case 2:
			htmlBodyTemp = htmlBodyTemp + "\n" + line
		}
	}
	if err := scanner.Err(); err != nil {
		log.Println("错误：无法读取资源文件: ", err)
	}
	// fmt.Println("htmlTemp", htmlTemp)
	// fmt.Println("htmlBodyTemp", htmlBodyTemp)
}

func genHTML() string {
	var html string = strings.Replace(htmlTemp, "!TITLE!", pageTitle, 1)
	var bodys []string = []string{}
	for _, v := range fontPathList {
		if len(v.Name) == 0 {
			continue
		}
		var body string = htmlBodyTemp
		body = strings.Replace(body, "!MD5!", v.MD5, -1)
		body = strings.Replace(body, "!NAME!", v.Name, -1)
		body = strings.Replace(body, "!FAMILY!", v.Family, -1)
		body = strings.Replace(body, "!STYLE!", v.SubfamilyName, -1)
		body = strings.Replace(body, "!FILE!", v.FontPath, -1)
		body = strings.Replace(body, "!VERSION!", strings.Replace(v.Version, "Version ", "", 1), -1)
		if v.Monospaced > 0 {
			body = strings.Replace(body, "!MONOSPACED!", "[字母等宽 "+strconv.Itoa(v.Monospaced)+"]", -1)
			body = strings.Replace(body, "!MONOSPACEDC!", "monospaced", -1)
		} else {
			body = strings.Replace(body, "!MONOSPACED!", "", -1)
			body = strings.Replace(body, "!MONOSPACEDC!", "", -1)
		}
		if v.MonospacedZH > -2 {
			if v.MonospacedZH > 0 {
				body = strings.Replace(body, "!ZH!", "[中文等宽 "+strconv.Itoa(v.MonospacedZH)+"]", -1)
				body = strings.Replace(body, "!ZHC!", " zh_monospaced", -1)
			} else {
				body = strings.Replace(body, "!ZH!", "[中文]", -1)
				body = strings.Replace(body, "!ZHC!", " zh", -1)
			}
		} else {
			body = strings.Replace(body, "!ZH!", "", -1)
			body = strings.Replace(body, "!ZHC!", "", -1)
		}
		bodys = append(bodys, body)
	}
	var body = strings.Join(bodys, "\n")
	return strings.Replace(html, "<!-- body -->", body, 1)
}

func saveHTML() {
	var html string = genHTML()
	// fmt.Println("html", html)
	var htmlPath string = outDir + "/index.html"
	os.WriteFile(htmlPath, []byte(html), 0644)
}
