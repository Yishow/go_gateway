package algorithms

import (
	"math"
	"testing"
)

// =============================================================================
// 任務 1.1: FFT 演算法測試
// =============================================================================

func TestFFT_DCOffset(t *testing.T) {
	// 輸入常數訊號 [1, 1, 1, 1]
	input := []complex128{1, 1, 1, 1}
	output := FFT(input)

	// DC 分量 (0Hz) 應該是 4 (所有值的總和)
	dcMagnitude := cmplxAbs(output[0])
	if math.Abs(dcMagnitude-4) > 0.001 {
		t.Errorf("DC 分量預期 4，實際 %.4f", dcMagnitude)
	}

	// 其他頻率分量應該接近 0
	for i := 1; i < len(output); i++ {
		mag := cmplxAbs(output[i])
		if mag > 0.001 {
			t.Errorf("頻率 %d 應該接近 0，實際 %.4f", i, mag)
		}
	}
}

func TestFFT_SineWave(t *testing.T) {
	// 產生 8 個樣本的正弦波 (頻率 = 1 個週期)
	n := 8
	input := make([]complex128, n)
	for i := 0; i < n; i++ {
		input[i] = complex(math.Sin(2*math.Pi*float64(i)/float64(n)), 0)
	}

	output := FFT(input)

	// 應該在 bin 1 (和對稱的 bin 7) 有峰值
	mag1 := cmplxAbs(output[1])
	mag0 := cmplxAbs(output[0])

	if mag1 < 3 {
		t.Errorf("Bin 1 應該有峰值 (約 4)，實際 %.4f", mag1)
	}
	if mag0 > 0.001 {
		t.Errorf("DC 分量應該接近 0，實際 %.4f", mag0)
	}
}

func TestFFT_PowerOfTwo(t *testing.T) {
	// 測試非 2 的冪次應該被處理
	input := []complex128{1, 2, 3, 4, 5, 6, 7, 8}
	output := FFT(input)

	if len(output) != len(input) {
		t.Errorf("輸出長度預期 %d，實際 %d", len(input), len(output))
	}
}

func TestFFT_Inverse(t *testing.T) {
	// 測試 IFFT 可以還原原始訊號
	original := []complex128{1, 2, 3, 4, 5, 6, 7, 8}
	spectrum := FFT(original)
	recovered := IFFT(spectrum)

	for i, v := range recovered {
		if math.Abs(real(v)-real(original[i])) > 0.001 {
			t.Errorf("還原失敗，索引 %d: 預期 %.4f，實際 %.4f", i, real(original[i]), real(v))
		}
	}
}

// cmplxAbs 計算複數絕對值
func cmplxAbs(c complex128) float64 {
	return math.Sqrt(real(c)*real(c) + imag(c)*imag(c))
}

// =============================================================================
// 任務 1.3: PID 控制器測試
// =============================================================================

func TestPID_ProportionalOnly(t *testing.T) {
	pid := NewPID(1.0, 0, 0) // Kp=1, Ki=0, Kd=0
	pid.SetTarget(100)

	// 輸入 0，誤差 = 100
	output := pid.Update(0, 0.1)

	if output != 100 {
		t.Errorf("P 控制：輸出預期 100，實際 %.4f", output)
	}
}

func TestPID_IntegralAccumulation(t *testing.T) {
	pid := NewPID(0, 1.0, 0) // Kp=0, Ki=1, Kd=0
	pid.SetTarget(100)

	// 多次更新，積分項應該累積
	for i := 0; i < 10; i++ {
		pid.Update(0, 0.1) // 每次誤差 100，dt=0.1，積分累積 100*0.1=10
	}

	output := pid.Update(0, 0.1)
	// 積分累積約 1100 (11 次 * 100 * 0.1)
	if output < 100 {
		t.Errorf("I 控制：輸出應該持續累積，實際 %.4f", output)
	}
}

func TestPID_DerivativeResponse(t *testing.T) {
	pid := NewPID(0, 0, 1.0) // Kp=0, Ki=0, Kd=1
	pid.SetTarget(100)

	// 第一次更新
	pid.Update(0, 0.1)
	// 第二次更新，輸入變化
	output := pid.Update(50, 0.1)

	// D 控制回應輸入變化率
	// 誤差從 100 變到 50，變化率 = -500/s
	if output >= 0 {
		t.Errorf("D 控制：誤差減少時輸出應為負，實際 %.4f", output)
	}
}

func TestPID_OutputLimits(t *testing.T) {
	pid := NewPID(10, 0, 0) // 高增益
	pid.SetTarget(1000)
	pid.SetOutputLimits(-100, 100)

	output := pid.Update(0, 0.1)

	if output > 100 || output < -100 {
		t.Errorf("輸出應該被限制在 [-100, 100]，實際 %.4f", output)
	}
}

func TestPID_Reset(t *testing.T) {
	pid := NewPID(1, 1, 0)
	pid.SetTarget(100)

	// 累積積分
	for i := 0; i < 10; i++ {
		pid.Update(0, 0.1)
	}

	pid.Reset()
	output := pid.Update(0, 0.1)

	// Reset 後，積分應該被清除
	// P 項：1 * 100 = 100
	// I 項：1 * 100 * 0.1 = 10 (第一次更新會累積)
	// 總計：110
	expectedMin := 100.0
	expectedMax := 115.0
	if output < expectedMin || output > expectedMax {
		t.Errorf("Reset 後輸出預期在 [%.0f, %.0f]，實際 %.4f", expectedMin, expectedMax, output)
	}
}
