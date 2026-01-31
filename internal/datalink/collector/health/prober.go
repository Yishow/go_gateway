package health

import (
	"context"
	"sync"
	"time"
)

// =============================================================================
// 背景探測器
// =============================================================================

// ProbeFunc 探測函數類型
type ProbeFunc func(ctx context.Context, deviceID string) error

// ProberConfig 探測器配置
type ProberConfig struct {
	// ProbeInterval 探測間隔 (Dead 狀態下的探測頻率)
	ProbeInterval time.Duration

	// ProbeTimeout 探測超時
	ProbeTimeout time.Duration
}

// DefaultProberConfig 預設探測器配置
func DefaultProberConfig() ProberConfig {
	return ProberConfig{
		ProbeInterval: 30 * time.Second,
		ProbeTimeout:  10 * time.Second,
	}
}

// BackgroundProber 背景探測器
// 負責對 Dead 狀態的設備發送輕量級心跳包
type BackgroundProber struct {
	mu sync.RWMutex

	config   ProberConfig
	breakers map[string]*CircuitBreaker
	probeFunc ProbeFunc

	running bool
	stopCh  chan struct{}
	wg      sync.WaitGroup
}

// NewBackgroundProber 建立新的背景探測器
func NewBackgroundProber(config ProberConfig, probeFunc ProbeFunc) *BackgroundProber {
	return &BackgroundProber{
		config:   config,
		breakers: make(map[string]*CircuitBreaker),
		probeFunc: probeFunc,
	}
}

// RegisterDevice 註冊設備熔斷器
func (p *BackgroundProber) RegisterDevice(deviceID string, breaker *CircuitBreaker) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.breakers[deviceID] = breaker
}

// UnregisterDevice 取消註冊設備
func (p *BackgroundProber) UnregisterDevice(deviceID string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	delete(p.breakers, deviceID)
}

// Start 啟動背景探測
func (p *BackgroundProber) Start() {
	p.mu.Lock()
	if p.running {
		p.mu.Unlock()
		return
	}
	p.running = true
	p.stopCh = make(chan struct{})
	p.mu.Unlock()

	p.wg.Add(1)
	go p.runProbeLoop()
}

// Stop 停止背景探測
func (p *BackgroundProber) Stop() {
	p.mu.Lock()
	if !p.running {
		p.mu.Unlock()
		return
	}
	p.running = false
	close(p.stopCh)
	p.mu.Unlock()

	p.wg.Wait()
}

// runProbeLoop 運行探測循環
func (p *BackgroundProber) runProbeLoop() {
	defer p.wg.Done()

	ticker := time.NewTicker(p.config.ProbeInterval)
	defer ticker.Stop()

	for {
		select {
		case <-p.stopCh:
			return
		case <-ticker.C:
			p.probeDeadDevices()
		}
	}
}

// probeDeadDevices 探測所有 Dead 狀態的設備
func (p *BackgroundProber) probeDeadDevices() {
	p.mu.RLock()
	deadDevices := make([]string, 0)
	for deviceID, breaker := range p.breakers {
		if breaker.State() == StateDead {
			deadDevices = append(deadDevices, deviceID)
		}
	}
	p.mu.RUnlock()

	// 並行探測
	var wg sync.WaitGroup
	for _, deviceID := range deadDevices {
		wg.Add(1)
		go func(devID string) {
			defer wg.Done()
			p.probeDevice(devID)
		}(deviceID)
	}
	wg.Wait()
}

// probeDevice 探測單一設備
func (p *BackgroundProber) probeDevice(deviceID string) {
	p.mu.RLock()
	breaker, exists := p.breakers[deviceID]
	p.mu.RUnlock()

	if !exists {
		return
	}

	// 檢查熔斷器是否允許探測
	if !breaker.AllowRequest() {
		return
	}

	// 執行探測
	ctx, cancel := context.WithTimeout(context.Background(), p.config.ProbeTimeout)
	defer cancel()

	err := p.probeFunc(ctx, deviceID)

	// 報告結果
	breaker.ReportResult(err)
}

// GetDeadDevices 取得所有 Dead 狀態的設備
func (p *BackgroundProber) GetDeadDevices() []string {
	p.mu.RLock()
	defer p.mu.RUnlock()

	deadDevices := make([]string, 0)
	for deviceID, breaker := range p.breakers {
		if breaker.State() == StateDead {
			deadDevices = append(deadDevices, deviceID)
		}
	}
	return deadDevices
}
