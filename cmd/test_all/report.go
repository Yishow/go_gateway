package main

import (
	"fmt"
	"time"
)

// addResult 添加測試結果
func (tr *TestResults) addResult(result TestResult) {
	tr.Results = append(tr.Results, result)
	tr.Total++
	if result.Success {
		tr.Passed++
	} else {
		tr.Failed++
	}
}

// printReport 輸出測試報告
func printReport(results *TestResults) {
	fmt.Println()
	fmt.Println("=========================================")
	fmt.Println("  測試報告")
	fmt.Println("=========================================")
	fmt.Println()

	for _, result := range results.Results {
		status := "✅ 通過"
		if !result.Success {
			status = "❌ 失敗"
		}
		fmt.Printf("%s [%s %s] %s (耗時: %v)\n",
			status,
			result.Protocol,
			result.Mode,
			result.Message,
			result.Duration.Round(time.Millisecond),
		)
	}

	fmt.Println()
	fmt.Println("=========================================")
	fmt.Printf("總計: %d  通過: %d  失敗: %d\n", results.Total, results.Passed, results.Failed)
	fmt.Println("=========================================")
}
