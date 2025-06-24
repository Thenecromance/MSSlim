package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"
)

//

func eraseModules() {
	log.Println("开始移除不必要的模块....")
	eraseBrowsePanel()

	patterns := []Pattern{
		{
			FileName:    path + "/" + "Module/ApplicantPanel.lua",
			Patterns:    "星标",
			OffsetStart: -2,
			OffsetEnd:   18,
		},
		{
			FileName:    path + "/" + "Module/ApplicantPanel.lua",
			Patterns:    "text = '@'",
			OffsetStart: -2,
			OffsetEnd:   18,
		},
		{
			FileName:    path + "/" + "Module/ApplicantPanel.lua",
			Patterns:    "操作",
			OffsetStart: -3,
			OffsetEnd:   5,
		},
		{
			FileName:    path + "/" + "Expansion/LocomotiveIntroduce.lua",
			Patterns:    "火车头",
			OffsetStart: 0,
			OffsetEnd:   0,
		},
		{
			FileName:    path + "/" + "Module/BrowsePanel.lua",
			Patterns:    "火车头",
			OffsetStart: -2,
			OffsetEnd:   0,
		},
	}
	eraseByPatterns(patterns)
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
		if strings.Contains(lines[i], "星标") ||
			strings.Contains(lines[i], "'@'") ||
			strings.Contains(lines[i], "'操作'") {
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

func eraseByPatterns(patterns []Pattern) {
	for _, pattern := range patterns {
		buf, err := os.ReadFile(pattern.FileName)
		if err != nil {
			fmt.Println("读取文件失败:", err)
			continue
		}
		content := string(buf)
		lines := strings.Split(content, "\n")
		for i, line := range lines {
			if strings.Contains(line, pattern.Patterns) {
				if strings.HasPrefix(line, "--") {
					continue // already disabled
				}
				if pattern.OffsetStart == 0 && pattern.OffsetEnd == 0 {
					lines[i] = "--" + lines[i] // disable the line
				} else {
					for idx := pattern.OffsetStart + i; idx <= i+pattern.OffsetEnd; idx++ {
						lines[idx] = "--" + lines[idx]
					}
				}

			}
		}

		os.WriteFile(pattern.FileName, []byte(strings.Join(lines, "\n")), 0644)
	}

}
