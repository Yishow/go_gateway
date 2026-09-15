package simulation

import (
	"math"
	"sync"
	"time"

	"go-gateway/internal/virtual/memory"
)

// =============================================================================
// 模擬規則類型
// =============================================================================

// RuleType 規則類型
type RuleType int

const (
	// RuleAutoIncrement 自動遞增
	RuleAutoIncrement RuleType = iota
	// RuleSineWave 正弦波
	RuleSineWave
	// RuleRandom 隨機值
	RuleRandom
	// RuleToggle 開關切換
	RuleToggle
)

// Rule 模擬規則
type Rule struct {
	ID        string
	RuleType  RuleType
	Offset    int // 記憶體偏移量
	Interval  time.Duration
	Increment int           // 用於 AutoIncrement
	MinValue  int           // 用於 Random
	MaxValue  int           // 用於 Random
	Amplitude uint16        // 用於 SineWave
	Period    time.Duration // 用於 SineWave
}

// =============================================================================
// SimulationEngine 模擬引擎
// =============================================================================

// SimulationEngine 虛擬設備模擬引擎
type SimulationEngine struct { //nolint:revive // Preserve the exported Go name and its existing callers during lint maintenance.
	mu      sync.RWMutex
	bank    *memory.MemoryBank
	rules   map[string]*Rule
	tickers map[string]*time.Ticker
	done    chan struct{}
	running bool
	wg      sync.WaitGroup
}

// NewSimulationEngine 建立新的模擬引擎
func NewSimulationEngine(bank *memory.MemoryBank) *SimulationEngine {
	return &SimulationEngine{
		bank:    bank,
		rules:   make(map[string]*Rule),
		tickers: make(map[string]*time.Ticker),
		done:    make(chan struct{}),
	}
}

// AddRule 新增模擬規則
func (e *SimulationEngine) AddRule(rule *Rule) {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.rules[rule.ID] = rule
}

// RemoveRule 移除模擬規則
func (e *SimulationEngine) RemoveRule(id string) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if ticker, exists := e.tickers[id]; exists {
		ticker.Stop()
		delete(e.tickers, id)
	}
	delete(e.rules, id)
}

// Start 啟動模擬引擎
func (e *SimulationEngine) Start() {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.running {
		return
	}

	e.running = true
	e.done = make(chan struct{})

	for _, rule := range e.rules {
		e.startRule(rule)
	}
}

// Stop 停止模擬引擎
func (e *SimulationEngine) Stop() {
	e.mu.Lock()
	if !e.running {
		e.mu.Unlock()
		return
	}
	e.running = false
	close(e.done)

	for id, ticker := range e.tickers {
		ticker.Stop()
		delete(e.tickers, id)
	}
	e.mu.Unlock()

	e.wg.Wait()
}

// startRule 啟動單一規則
func (e *SimulationEngine) startRule(rule *Rule) {
	ticker := time.NewTicker(rule.Interval)
	e.tickers[rule.ID] = ticker

	e.wg.Add(1)
	go func(r *Rule) {
		defer e.wg.Done()

		startTime := time.Now()
		counter := 0

		for {
			select {
			case <-e.done:
				return
			case <-ticker.C:
				e.executeRule(r, counter, startTime)
				counter++
			}
		}
	}(rule)
}

// executeRule 執行規則
func (e *SimulationEngine) executeRule(rule *Rule, counter int, startTime time.Time) {
	switch rule.RuleType {
	case RuleAutoIncrement:
		e.executeAutoIncrement(rule, counter)
	case RuleSineWave:
		e.executeSineWave(rule, startTime)
	case RuleRandom:
		e.executeRandom(rule)
	case RuleToggle:
		e.executeToggle(rule, counter)
	}
}

// executeAutoIncrement 執行自動遞增
func (e *SimulationEngine) executeAutoIncrement(rule *Rule, counter int) {
	value := safeUint16(counter * rule.Increment)
	if err := e.bank.WriteWord(rule.Offset, value); err != nil {
		return
	}
}

// executeSineWave 執行正弦波
func (e *SimulationEngine) executeSineWave(rule *Rule, startTime time.Time) {
	elapsed := time.Since(startTime).Seconds()
	periodSec := rule.Period.Seconds()
	if periodSec == 0 {
		periodSec = 1
	}

	// 正弦波計算: amplitude * sin(2π * t / period)
	sin := sinApprox(2 * 3.14159265359 * elapsed / periodSec)
	value := uint16(float64(rule.Amplitude) * (sin + 1) / 2) // 正規化到 0 ~ amplitude
	if err := e.bank.WriteWord(rule.Offset, value); err != nil {
		return
	}
}

// executeRandom 執行隨機值
func (e *SimulationEngine) executeRandom(rule *Rule) {
	valueRange := rule.MaxValue - rule.MinValue
	if valueRange <= 0 {
		valueRange = 1
	}
	// 簡單的偽隨機
	seed := time.Now().UnixNano()
	value := safeUint16(rule.MinValue + int(seed%int64(valueRange)))
	if err := e.bank.WriteWord(rule.Offset, value); err != nil {
		return
	}
}

// executeToggle 執行開關切換
func (e *SimulationEngine) executeToggle(rule *Rule, counter int) {
	value := safeUint16(counter % 2)
	if err := e.bank.WriteWord(rule.Offset, value); err != nil {
		return
	}
}

func safeUint16(value int) uint16 {
	if value < 0 {
		return 0
	}
	if value > math.MaxUint16 {
		return math.MaxUint16
	}
	return uint16(value)
}

// sinApprox 正弦近似函數 (避免 math 包依賴)
func sinApprox(x float64) float64 {
	// 泰勒展開近似
	for x > 3.14159265359 {
		x -= 2 * 3.14159265359
	}
	for x < -3.14159265359 {
		x += 2 * 3.14159265359
	}
	// sin(x) ≈ x - x³/6 + x⁵/120
	x3 := x * x * x
	x5 := x3 * x * x
	return x - x3/6 + x5/120
}
