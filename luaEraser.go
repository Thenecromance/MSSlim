package main

import (
	"bytes"
	"log"
	"os"
	"strings"
)

//

func eraseModules() {
	log.Println("开始移除不必要的模块....")
	eraseApplicantPanel()
	eraseBrowsePanel()

}

func eraseApplicantPanel() {
	ApplicantPanelPath := path + "/" + "Module/ApplicantPanel.lua"

	buf, err := os.ReadFile(ApplicantPanelPath)
	if err != nil {
		panic("读取文件失败" + ApplicantPanelPath)
	}

	lines := strings.Split(string(buf), "\n")
	start := findStartLine(lines, "local APPLICANT_LIST_HEADER = {")
	if start == -1 {
		panic("没有找到开始行")
	}

	blocks := enumTables(lines, start)
	for _, block := range blocks {
		if block.needToDelete {
			for i := block.start; i <= block.end; i++ {
				if !strings.Contains(lines[i], "--") {
					lines[i] = "--" + lines[i]
				}

			}
		}
	}

	buf = []byte(strings.Join(lines, "\n"))

	os.WriteFile(ApplicantPanelPath, buf, 0644)
}

func eraseBrowsePanel() {
	BrowsePanelPath := path + "/" + "Module/BrowsePanel.lua"
	buf, err := os.ReadFile(BrowsePanelPath)
	if err != nil {
		panic("读取文件失败")
	}
	// insert \n between }, { to make it easier to parse
	buf = bytes.Replace(buf, []byte("}, {"), []byte("},\n{"), -1)
	buf = bytes.Replace(buf, []byte("}, -- {"), []byte("},\n -- {"), -1)

	lines := strings.Split(string(buf), "\n")
	start := findStartLine(lines, "ActivityList:InitHeader {")
	if start == -1 {
		panic("没有找到开始行")
	}
	blocks := enumTables(lines, start)
	for _, block := range blocks {
		if block.needToDelete {
			for i := block.start; i <= block.end; i++ {
				if !strings.Contains(lines[i], "--") {
					lines[i] = "--" + lines[i]
				}

			}
		}
	}

	buf = []byte(strings.Join(lines, "\n"))

	os.WriteFile(BrowsePanelPath, buf, 0644)
}

func findStartLine(lines []string, start string) int {
	for i, line := range lines {
		if strings.Contains(line, start) {
			return i
		}
	}
	return -1
}

type block struct {
	start        int
	end          int
	needToDelete bool
}

func enumTables(lines []string, start int) []block {
	result := make([]block, 0)
	startLine := start + 1
	endLine := start + 1
	depth := 0
	needToDelete := false

	for i := startLine; i < len(lines); i++ {
		if strings.Contains(lines[i], "星标") || strings.Contains(lines[i], "'@'") || strings.Contains(lines[i], "'操作'") {
			needToDelete = true
		}

		if strings.Contains(lines[i], "{") {

			startLine = i
			depth++
		} else if strings.Contains(lines[i], "},") || strings.HasSuffix(lines[i], "}") {

			endLine = i
			depth--
		} else {
			continue
		}

		{
			if depth == 0 {
				result = append(result, block{
					start:        startLine,
					end:          endLine,
					needToDelete: needToDelete,
				})

				needToDelete = false
			}

			if depth < 0 { // when first time triggered } means this table is over
				break
			}
		}

	}

	return result
}
