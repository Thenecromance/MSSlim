package main

import (
	"log"
	"os"
)

func blackListFolder() {
	blackList := []string{
		"163UI_Info",
		"NewBeeBox",
		//"HeiBox"
	}
	for _, folder := range blackList {
		folderPath := path + "/../" + folder
		if err := os.RemoveAll(folderPath); err != nil {
			log.Printf("Failed to remove folder %s: %v\n", folder, err)
		} else {
			log.Printf("Successfully removed folder: %s\n", folder)
		}
	}
}
