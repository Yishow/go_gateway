package algorithms

// =============================================================================
// PID 控制器實作
// =============================================================================

// PID 控制器
type PID struct {
	kp, ki, kd float64 // PID 增益

	target    float64 // 目標值
	integral  float64 // 積分累積
	prevError float64 // 前次誤差

	minOutput, maxOutput float64 // 輸出限制
	hasLimits            bool
}

// NewPID 建立新的 PID 控制器
func NewPID(kp, ki, kd float64) *PID {
	return &PID{
		kp: kp,
		ki: ki,
		kd: kd,
	}
}

// SetTarget 設定目標值
func (p *PID) SetTarget(target float64) {
	p.target = target
}

// SetGains 設定 PID 增益
func (p *PID) SetGains(kp, ki, kd float64) {
	p.kp = kp
	p.ki = ki
	p.kd = kd
}

// SetOutputLimits 設定輸出限制
func (p *PID) SetOutputLimits(minimum, maximum float64) {
	p.minOutput = minimum
	p.maxOutput = maximum
	p.hasLimits = true
}

// Update 更新 PID 輸出
// input: 當前輸入值
// dt: 時間間隔 (秒)
// 返回: 控制輸出
func (p *PID) Update(input, dt float64) float64 {
	if dt <= 0 {
		dt = 0.001 // 避免除以零
	}

	// 計算誤差
	errValue := p.target - input

	// P 項
	pTerm := p.kp * errValue

	// I 項
	p.integral += errValue * dt
	iTerm := p.ki * p.integral

	// D 項
	derivative := (errValue - p.prevError) / dt
	dTerm := p.kd * derivative
	p.prevError = errValue

	// 總輸出
	output := pTerm + iTerm + dTerm

	// 應用輸出限制
	if p.hasLimits {
		if output > p.maxOutput {
			output = p.maxOutput
			// Anti-windup: 限制積分累積
			p.integral -= errValue * dt
		} else if output < p.minOutput {
			output = p.minOutput
			p.integral -= errValue * dt
		}
	}

	return output
}

// Reset 重置控制器狀態
func (p *PID) Reset() {
	p.integral = 0
	p.prevError = 0
}

// GetState 取得當前狀態
func (p *PID) GetState() (integral, prevError float64) {
	return p.integral, p.prevError
}

// =============================================================================
// 進階 PID 變體
// =============================================================================

// PIDWithFilter 帶有 D 項濾波的 PID
type PIDWithFilter struct {
	PID
	filterCoeff float64 // 濾波係數 (0-1)
	filteredD   float64 // 濾波後的 D 項
}

// NewPIDWithFilter 建立帶濾波的 PID
func NewPIDWithFilter(kp, ki, kd, filterCoeff float64) *PIDWithFilter {
	return &PIDWithFilter{
		PID: PID{
			kp: kp,
			ki: ki,
			kd: kd,
		},
		filterCoeff: filterCoeff,
	}
}

// Update 更新 PID 輸出 (帶 D 項濾波)
func (p *PIDWithFilter) Update(input, dt float64) float64 {
	if dt <= 0 {
		dt = 0.001
	}

	errValue := p.target - input

	// P 項
	pTerm := p.kp * errValue

	// I 項
	p.integral += errValue * dt
	iTerm := p.ki * p.integral

	// D 項 (帶濾波)
	derivative := (errValue - p.prevError) / dt
	p.filteredD = p.filterCoeff*p.filteredD + (1-p.filterCoeff)*derivative
	dTerm := p.kd * p.filteredD
	p.prevError = errValue

	output := pTerm + iTerm + dTerm

	if p.hasLimits {
		if output > p.maxOutput {
			output = p.maxOutput
			p.integral -= errValue * dt
		} else if output < p.minOutput {
			output = p.minOutput
			p.integral -= errValue * dt
		}
	}

	return output
}
