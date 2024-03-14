package main

import (
	"os"
	"strconv"
	"strings"
)

var htmlTemp string = `
<!DOCTYPE html>
<html lang="zh-CN">

<head>
	<meta charset="UTF-8">
	<meta name="viewport"
		content="width=device-width, initial-scale=1.0, maximum-scale=1.0, user-scalable=0,minimal-ui:ios">
	<meta http-equiv="X-UA-Compatible" content="ie=edge">
	<title>!TITLE!</title>
</head>

<body>
!BODY!
</body>
</html>
`

var htmlBodyTemp string = `
<p class="fontInfo" id="!MD5!">
    <div>
        <strong class="fontName">!NAME!</strong>&nbsp;
        <span class="fontFamily">!FAMILY!</span>&nbsp;
        <span class="fontFamily">(!STYLE!)</span>&nbsp;
		<span class="fontVersion">!VERSION!</span>
		<code class="fontMonospaced!MONOSPACEDC!">!MONOSPACED!</code>&nbsp;
		<code class="fontChinese!ZHC!">!ZH!</code>
    </div>
    <div>
        <code class="fontFile">!FILE!</code>&nbsp;
        <code class="fontMD5">(!MD5!)</code>&nbsp;
    </div>
    <img src="!MD5!.png" alt="[!NAME!] !MD5!.png" title="!NAME! [!FAMILY!] (!STYLE!)" />
</p>
`

func genHTML() string {
	var html string = htmlTemp
	for _, v := range fontPathList {
		if len(v.Name) == 0 {
			continue
		}
		var body string = strings.Replace(htmlBodyTemp, "!TITLE!", pageTitle, 1)
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
				body = strings.Replace(body, "!ZHC!", "zh_monospaced", -1)
			} else {
				body = strings.Replace(body, "!ZH!", "[中文]", -1)
				body = strings.Replace(body, "!ZHC!", "zh", -1)
			}
		} else {
			body = strings.Replace(body, "!ZH!", "", -1)
			body = strings.Replace(body, "!ZHC!", "", -1)
		}
		html = strings.Replace(html, "!BODY!", body+"!BODY!", -1)
	}
	return strings.Replace(html, "!BODY!", "", 1)
}

func saveHTML() {
	var html string = genHTML()
	var htmlPath string = outDir + "/index.html"
	os.WriteFile(htmlPath, []byte(html), 0644)
}
