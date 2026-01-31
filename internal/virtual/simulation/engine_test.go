package simulation

import (
	"testing"
	"time"

	"go-gateway/internal/virtual/memory"
)

// =============================================================================
// 任務 3.1: SimulationEngine 測試
// =============================================================================

func TestSimulationEngine_AutoIncrement(t *testing.T) {
	bank := memory.NewMemoryBank(1024)
	engine := NewSimulationEngine(bank)

	// 新增自動遞增規則: 每 100ms 遞增 1
	rule := &Rule{
		ID:        "test-inc",
		RuleType:  RuleAutoIncrement,
		Offset:    0,           // D0
		Interval:  100 * time.Millisecond,
		Increment: 1,
	}
	engine.AddRule(rule)

	// 啟動引擎
	engine.Start()
	defer engine.Stop()

	// 等待 2 秒
	time.Sleep(2 * time.Second)

	// 讀取值，應該已經遞增多次
	val, _ := bank.ReadWord(0)
	
	// 2 秒內應該執行約 20 次
	if val < 15 || val > 25 {
		t.Errorf("預期約 20 次遞增，實際值: %d", val)
	}
}

func TestSimulationEngine_Toggle(t *testing.T) {
	bank := memory.NewMemoryBank(1024)
	engine := NewSimulationEngine(bank)

	rule := &Rule{
		ID:       "test-toggle",
		RuleType: RuleToggle,
		Offset:   10,
		Interval: 50 * time.Millisecond,
	}
	engine.AddRule(rule)

	engine.Start()
	defer engine.Stop()

	// 收集一系列值
	values := make([]uint16, 0)
	for i := 0; i < 5; i++ {
		time.Sleep(60 * time.Millisecond)
		val, _ := bank.ReadWord(10)
		values = append(values, val)
	}

	// 檢查是否有切換
	hasZero := false
	hasOne := false
	for _, v := range values {
		if v == 0 {
			hasZero = true
		}
		if v == 1 {
			hasOne = true
		}
	}

	if !hasZero || !hasOne {
		t.Errorf("預期有 0 和 1 的切換，實際值: %v", values)
	}
}

func TestSimulationEngine_StopStart(t *testing.T) {
	bank := memory.NewMemoryBank(1024)
	engine := NewSimulationEngine(bank)

	rule := &Rule{
		ID:        "test-stop",
		RuleType:  RuleAutoIncrement,
		Offset:    20,
		Interval:  50 * time.Millisecond,
		Increment: 10,
	}
	engine.AddRule(rule)

	// 啟動
	engine.Start()
	time.Sleep(200 * time.Millisecond)
	engine.Stop()

	// 記錄停止時的值
	valAfterStop, _ := bank.ReadWord(20)

	// 等待一段時間
	time.Sleep(200 * time.Millisecond)

	// 值應該不再變化
	valLater, _ := bank.ReadWord(20)
	
	if valLater != valAfterStop {
		t.Errorf("停止後值不應變化: 停止時 %d，之後 %d", valAfterStop, valLater)
	}
}

func TestSimulationEngine_MultipleRules(t *testing.T) {
	bank := memory.NewMemoryBank(1024)
	engine := NewSimulationEngine(bank)

	// 新增多個規則
	engine.AddRule(&Rule{
		ID:        "rule1",
		RuleType:  RuleAutoIncrement,
		Offset:    0,
		Interval:  100 * time.Millisecond,
		Increment: 1,
	})

	engine.AddRule(&Rule{
		ID:       "rule2",
		RuleType: RuleToggle,
		Offset:   10,
		Interval: 100 * time.Millisecond,
	})

	engine.Start()
	time.Sleep(500 * time.Millisecond)
	engine.Stop()

	// 檢查兩個偏移量都有值
	val1, _ := bank.ReadWord(0)
	val2, _ := bank.ReadWord(10)

	if val1 == 0 && val2 == 0 {
		t.Error("至少一個規則應該產生非零值")
	}
}
