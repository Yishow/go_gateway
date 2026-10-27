package dbtarget

import (
	"errors"
	"strings"

	"go-gateway/internal/datalink/schema"
)

const sqlServerTLSVerificationFailed = "tls_verification_failed"

func classifyConnectorError(err error) schema.DatabaseConnectorStatus {
	var numbered interface{ SQLErrorNumber() int32 }
	if errors.As(err, &numbered) && numbered.SQLErrorNumber() == 18456 {
		return schema.DatabaseConnectorStatusAuthFailed
	}
	return classifyOtherConnectorError(err)
}

func connectorErrorMessage(kind schema.DatabaseConnectorKind, err error) string {
	if kind != schema.DatabaseConnectorKindSQLServer {
		return otherConnectorErrorMessage(kind, err)
	}
	switch sqlServerConnectionFailureReason(err) {
	case sqlServerTLSVerificationFailed:
		return "SQL Server TLS 憑證驗證失敗，請檢查已儲存的 CA 與伺服器名稱。"
	case inspectionReasonAuthenticationFailed:
		return "SQL Server 驗證失敗，請檢查已儲存的 SQL 帳密。"
	default:
		return "SQL Server 連線失敗，請檢查目的地連線設定及服務狀態。"
	}
}

func sqlServerConnectionFailureReason(err error) string {
	if err == nil {
		return inspectionReasonConnectionFailed
	}
	text := strings.ToLower(err.Error())
	if strings.Contains(text, "certificate") || strings.Contains(text, "x509") || strings.Contains(text, "tls") {
		return sqlServerTLSVerificationFailed
	}
	if classifyConnectorError(err) == schema.DatabaseConnectorStatusAuthFailed {
		return inspectionReasonAuthenticationFailed
	}
	return inspectionReasonConnectionFailed
}
