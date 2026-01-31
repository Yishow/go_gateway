package edge

import (
	"context"
	"time"

	"gopkg.in/yaml.v3"
)

// =============================================================================
// Processor 介面
// =============================================================================

// Processor 資料處理器介面
type Processor interface {
	// Name 返回處理器名稱
	Name() string
	// Process 處理數據
	Process(ctx context.Context, input []float64) ([]float64, error)
}

// =============================================================================
// Pipeline 資料處理管道
// =============================================================================

// Pipeline 資料處理管道
type Pipeline struct {
	processors []Processor
	timeout    time.Duration
}

// NewPipeline 建立新的處理管道
func NewPipeline() *Pipeline {
	return &Pipeline{
		processors: make([]Processor, 0),
	}
}

// AddProcessor 新增處理器
func (p *Pipeline) AddProcessor(proc Processor) {
	p.processors = append(p.processors, proc)
}

// SetTimeout 設定超時時間
func (p *Pipeline) SetTimeout(timeout time.Duration) {
	p.timeout = timeout
}

// Run 執行處理管道
func (p *Pipeline) Run(ctx context.Context, input []float64) ([]float64, error) {
	// 如果設定了超時，建立新的 context
	if p.timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, p.timeout)
		defer cancel()
	}

	// 複製輸入以避免修改原始數據
	data := make([]float64, len(input))
	copy(data, input)

	// 依序執行處理器
	for _, proc := range p.processors {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		var err error
		data, err = proc.Process(ctx, data)
		if err != nil {
			return nil, err
		}
	}

	return data, nil
}

// =============================================================================
// 配置解析
// =============================================================================

// PipelineConfig 管道配置
type PipelineConfig struct {
	Processors []ProcessorConfig `yaml:"processors"`
	Timeout    string            `yaml:"timeout,omitempty"`
}

// ProcessorConfig 處理器配置
type ProcessorConfig struct {
	Type   string                 `yaml:"type"`
	Params map[string]interface{} `yaml:"params,omitempty"`
}

// ParsePipelineConfig 解析 YAML 配置
func ParsePipelineConfig(data []byte) (*PipelineConfig, error) {
	var config PipelineConfig
	err := yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

// =============================================================================
// ProcessorRegistry 處理器註冊表
// =============================================================================

// ProcessorFactory 處理器工廠函數
type ProcessorFactory func(params map[string]interface{}) (Processor, error)

// ProcessorRegistry 處理器註冊表
type ProcessorRegistry struct {
	factories map[string]ProcessorFactory
}

// NewProcessorRegistry 建立新的註冊表
func NewProcessorRegistry() *ProcessorRegistry {
	return &ProcessorRegistry{
		factories: make(map[string]ProcessorFactory),
	}
}

// Register 註冊處理器工廠
func (r *ProcessorRegistry) Register(name string, factory ProcessorFactory) {
	r.factories[name] = factory
}

// Create 建立處理器
func (r *ProcessorRegistry) Create(name string, params map[string]interface{}) (Processor, error) {
	factory, exists := r.factories[name]
	if !exists {
		return nil, nil // 或返回錯誤
	}
	return factory(params)
}

// BuildPipeline 從配置建立管道
func (r *ProcessorRegistry) BuildPipeline(config *PipelineConfig) (*Pipeline, error) {
	pipeline := NewPipeline()

	// 設定超時
	if config.Timeout != "" {
		duration, err := time.ParseDuration(config.Timeout)
		if err == nil {
			pipeline.SetTimeout(duration)
		}
	}

	// 建立處理器
	for _, pc := range config.Processors {
		proc, err := r.Create(pc.Type, pc.Params)
		if err != nil {
			return nil, err
		}
		if proc != nil {
			pipeline.AddProcessor(proc)
		}
	}

	return pipeline, nil
}
