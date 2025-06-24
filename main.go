package main

import (
	"github.com/spf13/cobra"
	"log"
)

const Version = "1.0.1"

var path string
var warning = "注意：该exe仅为个人使用， 不对任何人负责，使用前请备份集合石插件的文件\n" +
	"正常游戏时请使用集合石插件的最新版本\n"
var cleanerCmd = &cobra.Command{
	Short: "清理集合石",
	Long:  "去除集合石中不必要的内容",
	Run: func(cmd *cobra.Command, args []string) {
		log.Println(warning)
		log.Println("集合石清理工具版本:", Version)
		log.Println("开始清理集合石....")
		defer func() {
			log.Println("清理完成")
			log.Println("如果游戏在线，请使用/rl命令重载集合石插件")
		}()
		eraseLoadXml()
		eraseModules()
		blackListFolder()

	},
}

func init() {
	cleanerCmd.Flags().StringVarP(&path, "path", "p", "", "指定集合石插件的路径")
	cleanerCmd.MarkFlagRequired("path")

}

func main() {
	if cleanerCmd.Execute() != nil {
		log.Fatal("Error executing command")
	}
}
