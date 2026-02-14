package edge

import (
	"context"
	"errors"
	"testing"
	"time"
)

// =============================================================================
// 任務 2.1: Pipeline 處理測試
// =============================================================================

// MockAddOneProcessor 測試用處理器，將每個值 +1
type MockAddOneProcessor struct{}

func (p *MockAddOneProcessor) Name() string { return "add_one" }

func (p *MockAddOneProcessor) Process(ctx context.Context, input []float64) ([]float64, error) {
	result := make([]float64, len(input))
	for i, v := range input {
		result[i] = v + 1
	}
	return result, nil
}

// MockMultiplyProcessor 測試用處理器，將每個值乘以倍數
type MockMultiplyProcessor struct {
	factor float64
}

func (p *MockMultiplyProcessor) Name() string { return "multiply" }

func (p *MockMultiplyProcessor) Process(ctx context.Context, input []float64) ([]float64, error) {
	result := make([]float64, len(input))
	for i, v := range input {
		result[i] = v * p.factor
	}
	return result, nil
}

func TestPipeline_SingleProcessor(t *testing.T) {
	pipeline := NewPipeline()
	pipeline.AddProcessor(&MockAddOneProcessor{})

	input := []float64{1, 2, 3}
	result, err := pipeline.Run(context.Background(), input)

	if err != nil {
		t.Fatalf("Pipeline 執行失敗: %v", err)
	}

	expected := []float64{2, 3, 4}
	for i, v := range result {
		if v != expected[i] {
			t.Errorf("索引 %d: 預期 %.1f，實際 %.1f", i, expected[i], v)
		}
	}
}

func TestPipeline_MultipleProcessors(t *testing.T) {
	pipeline := NewPipeline()
	pipeline.AddProcessor(&MockAddOneProcessor{})            // +1
	pipeline.AddProcessor(&MockMultiplyProcessor{factor: 2}) // *2

	input := []float64{1, 2, 3}
	result, err := pipeline.Run(context.Background(), input)

	if err != nil {
		t.Fatalf("Pipeline 執行失敗: %v", err)
	}

	// (1+1)*2=4, (2+1)*2=6, (3+1)*2=8
	expected := []float64{4, 6, 8}
	for i, v := range result {
		if v != expected[i] {
			t.Errorf("索引 %d: 預期 %.1f，實際 %.1f", i, expected[i], v)
		}
	}
}

func TestPipeline_EmptyInput(t *testing.T) {
	pipeline := NewPipeline()
	pipeline.AddProcessor(&MockAddOneProcessor{})

	input := []float64{}
	result, err := pipeline.Run(context.Background(), input)

	if err != nil {
		t.Fatalf("Pipeline 執行失敗: %v", err)
	}

	if len(result) != 0 {
		t.Errorf("空輸入應返回空結果")
	}
}

func TestPipeline_NoProcessors(t *testing.T) {
	pipeline := NewPipeline()

	input := []float64{1, 2, 3}
	result, err := pipeline.Run(context.Background(), input)

	if err != nil {
		t.Fatalf("Pipeline 執行失敗: %v", err)
	}

	// 無處理器時應返回原始輸入
	for i, v := range result {
		if v != input[i] {
			t.Errorf("無處理器應返回原始輸入")
		}
	}
}

// =============================================================================
// 任務 4.1: 超時丟棄機制測試
// =============================================================================

// MockSlowProcessor 模擬慢速處理器
type MockSlowProcessor struct {
	delay time.Duration
}

func (p *MockSlowProcessor) Name() string { return "slow" }

func (p *MockSlowProcessor) Process(ctx context.Context, input []float64) ([]float64, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-time.After(p.delay):
		return input, nil
	}
}

func TestPipeline_Timeout(t *testing.T) {
	pipeline := NewPipeline()
	pipeline.AddProcessor(&MockSlowProcessor{delay: 1 * time.Second})
	pipeline.SetTimeout(100 * time.Millisecond)

	input := []float64{1, 2, 3}
	_, err := pipeline.Run(context.Background(), input)

	if err == nil {
		t.Error("預期超時錯誤")
	}
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("預期 DeadlineExceeded 錯誤，實際: %v", err)
	}
}

func TestPipeline_TimeoutDoesNotBlock(t *testing.T) {
	pipeline := NewPipeline()
	pipeline.AddProcessor(&MockSlowProcessor{delay: 1 * time.Second})
	pipeline.SetTimeout(50 * time.Millisecond)

	start := time.Now()
	_, _ = pipeline.Run(context.Background(), []float64{1})
	elapsed := time.Since(start)

	// 應該在 100ms 內返回，而不是等待 1s
	if elapsed > 200*time.Millisecond {
		t.Errorf("超時應該快速返回，實際耗時: %v", elapsed)
	}
}

// =============================================================================
// 任務 3.1: 配置解析測試
// =============================================================================

func TestPipelineConfig_Parse(t *testing.T) {
	yaml := `
processors:
  - type: add_one
  - type: multiply
    params:
      factor: 2
`
	config, err := ParsePipelineConfig([]byte(yaml))
	if err != nil {
		t.Fatalf("解析配置失敗: %v", err)
	}

	if len(config.Processors) != 2 {
		t.Errorf("預期 2 個處理器，實際 %d", len(config.Processors))
	}

	if config.Processors[0].Type != "add_one" {
		t.Errorf("第一個處理器類型預期 add_one，實際 %s", config.Processors[0].Type)
	}

	if config.Processors[1].Type != "multiply" {
		t.Errorf("第二個處理器類型預期 multiply，實際 %s", config.Processors[1].Type)
	}
}
