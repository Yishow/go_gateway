// Package main provides the gateway CLI tool.
//
// Usage:
//
//	gateway [command]
//
// Available Commands:
//
//	validate    Validate configuration file
//	version     Show version information
//	help        Help about any command
package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var version = "0.1.0"

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "gateway",
	Short: "Go Gateway - 工業數據採集閘道",
	Long: `Go Gateway 是一個工業級數據採集閘道系統，
用於從 PLC (Programmable Logic Controller) 採集數據並映射至時序資料庫。

支援協議：Modbus RTU/TCP、FATEK FBs、MC Protocol 3E、MQTT`,
}

func init() {
	rootCmd.AddCommand(validateCmd)
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "顯示版本資訊",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Go Gateway v%s\n", version)
	},
}
