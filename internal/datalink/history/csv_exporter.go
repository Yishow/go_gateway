package history

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"time"
)

func sanitizeCSVField(s string) string {
	if s == "" {
		return ""
	}
	// 若是純數值（包含負數如 -12.3），不跳脫
	if _, err := strconv.ParseFloat(s, 64); err == nil {
		return s
	}

	// 若以可能觸發公式的符號開頭，加上單引號前綴
	firstChar := s[0]
	if firstChar == '=' || firstChar == '+' || firstChar == '-' || firstChar == '@' || firstChar == '\t' || firstChar == '\r' {
		return "'" + s
	}
	return s
}

// ExportCSV 將歷史資料點以 RFC4180 標準與 Formula Injection 防護串流寫入 CSV。
func ExportCSV(w io.Writer, points []HistoryPoint) error {
	writer := csv.NewWriter(w)
	defer writer.Flush()

	header := []string{
		"observed_at",
		"measurement_id",
		"value_numeric",
		"value_string",
		"time_weighted_mean",
		"sampled_min",
		"sampled_max",
		"usage_delta",
		"quality",
		"coverage_ratio",
		"is_estimated",
		"is_provisional",
	}

	if err := writer.Write(header); err != nil {
		return err
	}

	for _, p := range points {
		valNum := ""
		if p.ValueNumeric != nil {
			valNum = fmt.Sprintf("%.6f", *p.ValueNumeric)
		}

		valStr := ""
		if p.ValueString != nil {
			valStr = sanitizeCSVField(*p.ValueString)
		}

		mean := ""
		if p.TimeWeightedMean != nil {
			mean = fmt.Sprintf("%.6f", *p.TimeWeightedMean)
		}

		minVal := ""
		if p.SampledMin != nil {
			minVal = fmt.Sprintf("%.6f", *p.SampledMin)
		}

		maxVal := ""
		if p.SampledMax != nil {
			maxVal = fmt.Sprintf("%.6f", *p.SampledMax)
		}

		usage := ""
		if p.UsageDelta != nil {
			usage = fmt.Sprintf("%.6f", *p.UsageDelta)
		}

		row := []string{
			p.ObservedAt.Format(time.RFC3339),
			p.MeasurementID,
			valNum,
			valStr,
			mean,
			minVal,
			maxVal,
			usage,
			p.Quality,
			fmt.Sprintf("%.4f", p.CoverageRatio),
			strconv.FormatBool(p.IsEstimated),
			strconv.FormatBool(p.IsProvisional),
		}

		if err := writer.Write(row); err != nil {
			return err
		}
	}

	return nil
}
