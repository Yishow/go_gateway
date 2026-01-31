package main

import (
	"encoding/json"
	"fmt"
	"os"

	"go-gateway/internal/config"

	"github.com/spf13/cobra"
)

var (
	configFile string
	jsonOutput bool
)

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "驗證配置文件",
	Long: `驗證配置文件的格式和內容是否正確。

範例:
  gateway validate                    # 驗證預設配置
  gateway validate -c config.yaml     # 驗證指定配置文件
  gateway validate --json             # 以 JSON 格式輸出結果`,
	RunE: runValidate,
}

func init() {
	validateCmd.Flags().StringVarP(&configFile, "config", "c", "", "配置文件路徑")
	validateCmd.Flags().BoolVar(&jsonOutput, "json", false, "以 JSON 格式輸出")
}

// ValidationOutput 驗證輸出結構
type ValidationOutput struct {
	Valid   bool                      `json:"valid"`
	Errors  []config.ValidationError  `json:"errors,omitempty"`
	Summary string                    `json:"summary"`
}

func runValidate(cmd *cobra.Command, args []string) error {
	// Load configuration
	var cfg *config.Config
	var err error

	if configFile != "" {
		cfg, err = config.LoadFromFile(configFile)
	} else {
		cfg, err = config.Load()
	}

	if err != nil {
		return fmt.Errorf("無法載入配置: %w", err)
	}

	// Validate configuration
	result := config.ValidateConfig(cfg)

	// Output results
	if jsonOutput {
		return outputJSON(result)
	}

	return outputText(result)
}

func outputJSON(result *config.ValidationResult) error {
	output := ValidationOutput{
		Valid:  result.Valid,
		Errors: result.Errors,
	}

	if result.Valid {
		output.Summary = "配置驗證通過"
	} else {
		output.Summary = fmt.Sprintf("配置驗證失敗，共 %d 個錯誤", len(result.Errors))
	}

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(output)
}

func outputText(result *config.ValidationResult) error {
	if result.Valid {
		fmt.Println("✅ 配置驗證通過")
		return nil
	}

	fmt.Println("❌ 配置驗證失敗")
	fmt.Println()

	for i, err := range result.Errors {
		fmt.Printf("  %d. [%s] %s\n", i+1, err.Field, err.Message)
	}

	fmt.Println()
	fmt.Printf("共 %d 個錯誤\n", len(result.Errors))

	return fmt.Errorf("配置驗證失敗")
}
