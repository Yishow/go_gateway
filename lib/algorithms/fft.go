package algorithms

import (
	"math"
	"math/cmplx"
)

// =============================================================================
// FFT (Fast Fourier Transform) 實作
// Cooley-Tukey 算法
// =============================================================================

// FFT 執行快速傅立葉變換
func FFT(x []complex128) []complex128 {
	n := len(x)
	if n <= 1 {
		return x
	}

	// 確保長度是 2 的冪次
	if n&(n-1) != 0 {
		// 補零到最近的 2 的冪次
		newN := nextPowerOfTwo(n)
		padded := make([]complex128, newN)
		copy(padded, x)
		x = padded
		n = newN
	}

	// 遞迴實作
	return fftRecursive(x)
}

// fftRecursive Cooley-Tukey 遞迴算法
func fftRecursive(x []complex128) []complex128 {
	n := len(x)
	if n <= 1 {
		return x
	}

	// 分離偶數和奇數索引
	even := make([]complex128, n/2)
	odd := make([]complex128, n/2)
	for i := 0; i < n/2; i++ {
		even[i] = x[2*i]
		odd[i] = x[2*i+1]
	}

	// 遞迴處理
	even = fftRecursive(even)
	odd = fftRecursive(odd)

	// 合併
	result := make([]complex128, n)
	for k := 0; k < n/2; k++ {
		// 旋轉因子 W_n^k = e^(-2πik/n)
		angle := -2.0 * math.Pi * float64(k) / float64(n)
		w := cmplx.Exp(complex(0, angle))
		t := w * odd[k]
		result[k] = even[k] + t
		result[k+n/2] = even[k] - t
	}

	return result
}

// IFFT 執行逆快速傅立葉變換
func IFFT(x []complex128) []complex128 {
	n := len(x)
	if n <= 1 {
		return x
	}

	// 取共軛
	conj := make([]complex128, n)
	for i := range x {
		conj[i] = cmplx.Conj(x[i])
	}

	// 執行 FFT
	result := FFT(conj)

	// 取共軛並除以 N
	for i := range result {
		result[i] = cmplx.Conj(result[i]) / complex(float64(n), 0)
	}

	return result
}

// nextPowerOfTwo 找到下一個 2 的冪次
func nextPowerOfTwo(n int) int {
	p := 1
	for p < n {
		p *= 2
	}
	return p
}

// =============================================================================
// 頻譜分析輔助函數
// =============================================================================

// Magnitude 計算頻譜幅度
func Magnitude(spectrum []complex128) []float64 {
	result := make([]float64, len(spectrum))
	for i, v := range spectrum {
		result[i] = cmplx.Abs(v)
	}
	return result
}

// Phase 計算頻譜相位
func Phase(spectrum []complex128) []float64 {
	result := make([]float64, len(spectrum))
	for i, v := range spectrum {
		result[i] = cmplx.Phase(v)
	}
	return result
}

// PowerSpectrum 計算功率頻譜
func PowerSpectrum(spectrum []complex128) []float64 {
	result := make([]float64, len(spectrum))
	for i, v := range spectrum {
		mag := cmplx.Abs(v)
		result[i] = mag * mag
	}
	return result
}

// DominantFrequency 找到主頻率 (以 bin 索引返回)
func DominantFrequency(spectrum []complex128) int {
	maxMag := 0.0
	maxIdx := 0

	// 只檢查前半部分 (對稱性)
	n := len(spectrum) / 2
	for i := 1; i < n; i++ { // 跳過 DC
		mag := cmplx.Abs(spectrum[i])
		if mag > maxMag {
			maxMag = mag
			maxIdx = i
		}
	}

	return maxIdx
}
