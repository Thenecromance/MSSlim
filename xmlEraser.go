package main

import (
	"log"
	"os"
	"strings"
)

func eraseLoadXml() {
	log.Println("开始移除不必要的Load.xml....")
	xmlPath := path + "/" + "Load.xml"
	buf, err := os.ReadFile(xmlPath)
	if err != nil {
		log.Println("读取文件失败:", err)
		return
	}
	/*defer os.WriteFile(xmlPath, buf, 0644)*/
	lines := strings.Split(string(buf), "\n")
	for i, line := range lines {
		if strings.Contains(line, "<!--") || strings.Contains(line, "-->") {
			continue // skip the line
		}

		if strings.Contains(line, "Data\\Load.xml") {
			lines[i] = "<!--" + line + "-->" // disable the line
		}
	}
	buf = []byte(strings.Join(lines, "\n"))
	err = os.WriteFile(xmlPath, buf, 0644)
	if err != nil {
		log.Println("写入文件失败:", err)
		return
	}

}
