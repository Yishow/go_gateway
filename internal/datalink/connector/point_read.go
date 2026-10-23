package connector

import "go-gateway/internal/datalink/schema"

// PointReadCount returns the protocol-specific quantity for one typed point.
// MC counts typed values, while the existing Modbus/FATEK adapters count words.
// Direct batch callers retain their own Count contract.
func PointReadCount(protocol schema.ProtocolType, dataType schema.DataType) int {
	if protocol == schema.ProtocolMC3E {
		return 1
	}
	return schema.RegisterCountForDataType(dataType)
}
