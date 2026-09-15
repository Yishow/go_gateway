package collector

import (
	"context"
	"time"

	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/schema"
)

// pollPoint 輪詢單一點位，返回錯誤用於熔斷器報告
func (s *Scheduler) pollPoint(conn *connector.ManagedConnection, pt pointInfo) error {
	req := connector.ReadRequest{
		Address:    pt.Address,
		Function:   pt.Function,
		DataType:   pt.DataType,
		DataFormat: pt.DataFormat,
	}

	req.Count = schema.RegisterCountForDataType(pt.DataType)

	var result connector.ReadResult
	var err error

	// 重試邏輯 (指數退避)
	for attempt := 0; attempt <= s.config.DefaultRetryCount; attempt++ {
		if attempt > 0 {
			// 指數退避: delay * 2^attempt
			backoff := s.config.DefaultRetryDelay * time.Duration(1<<uint(attempt))
			if backoff > 30*time.Second {
				backoff = 30 * time.Second // 最大 30 秒
			}
			time.Sleep(backoff)
		}

		result, err = conn.Read(context.Background(), req)
		if err == nil && result.Quality == schema.QualityGood {
			break
		}
	}

	// 發送結果
	cv := CollectedValue{
		PointID:   pt.ID,
		DeviceID:  pt.DeviceID,
		Value:     result.Value,
		RawBytes:  result.RawBytes,
		Timestamp: result.Timestamp,
		Quality:   result.Quality,
		Error:     result.Error,
	}

	if err != nil {
		cv.Quality = schema.QualityBad
		cv.Error = err.Error()
	}

	s.emitValue(cv)
	return err
}

// emitValue 發送收集到的值
func (s *Scheduler) emitValue(cv CollectedValue) {
	select {
	case s.valueChan <- cv:
	default:
		// 緩衝區已滿，嘗試丟棄舊值後再非阻塞寫入
		select {
		case <-s.valueChan:
		default:
		}
		select {
		case s.valueChan <- cv:
		default:
		}
	}
}
